package rb

import (
	"context"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"net/http"
	"strconv"

	"go.uber.org/zap"
)

const (
	// CSRFTokenLength is the token length
	CSRFTokenLength = 32

	// CSRFSessionFieldName configures part of the session that will be checked
	CSRFSessionFieldName = "_rb_csrf"

	// CSRFFormFieldName configures which part of the form body will be checked
	CSRFFormFieldName = "rb.csrf.token"

	// CSRFHeaderName configures which header will be checked for csrf token
	CSRFHeaderName = "X-RB-CSRF-Token"
)

var (
	// ErrBadCSRFToken is returned if the CSRF token in the request does not match
	// the token in the session, or is otherwise malformed.
	ErrBadCSRFToken = errors.New("CSRF token invalid")
)

// CSRFToken returns the id token from the context if it has any
func CSRFToken(ctx context.Context) (tok string) {
	tok, _ = ctx.Value(ctxKey("token")).(string)
	return
}

// WithCSRFToken sets the id token context value
func WithCSRFToken(ctx context.Context, tok string) context.Context {
	return context.WithValue(ctx, ctxKey("token"), tok)
}

// CSRFMiddleware protects all non-get methods from cross-site request forgery
// by comparing a secret token in the users's session with one provided in the
// request body or header
type CSRFMiddleware func(http.Handler) http.Handler

// NewCSRFMiddlware creates the actual middleware
func NewCSRFMiddlware(sc SessionCore, ec ErrorCore) CSRFMiddleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			// check if there is a valid token in the session
			realToken, ok := sc.Session(w, r).Get(CSRFSessionFieldName).([]byte)
			if !ok || len(realToken) != CSRFTokenLength {

				L(r).Debug("no valid CSRF token in session, generating new one",
					zap.Bool("present", ok), zap.Int("len", len(realToken)))

				// if not, generate new random bytes
				realToken = gen(CSRFTokenLength)

				// Save the new (real) token in the session store.
				sc.Session(w, r).Set(CSRFSessionFieldName, realToken)
			}

			// mask the token to protect again BREACH and encode using base64
			masked := mask(realToken)

			// store in the request context so it can be retrieved in handlers
			r = r.WithContext(WithCSRFToken(r.Context(), masked))

			// unsafe methods need inspection
			if _, ok := safeMethods[r.Method]; !ok {
				L(r).Debug("unsafe method, comparing CSRF tokens", zap.String("method", r.Method))

				// Retrieve the combined token (pad + masked) token and unmask it.
				requestToken := unmask(requestToken(r))

				// Compare with constant time
				if subtle.ConstantTimeCompare(requestToken, realToken) == 0 {
					L(r).Debug("failed CSRF token compare",
						zap.Int("len_req_token", len(requestToken)), zap.Int("len_real_token", len(realToken)))

					ec.HandleError(w, r, ErrBadCSRFToken)
					return
				}
			}

			// call the next handler
			next.ServeHTTP(w, r)

			// Set the Vary: Cookie header to protect clients from caching the response.
			w.Header().Add("Vary", "Cookie")
		})
	}
}

// requestToken returns the issued token (pad + masked token) from the HTTP POST
// body or HTTP header. It will return nil if the token fails to decode.
func requestToken(r *http.Request) []byte {
	// 1. Check the HTTP header first.
	issued := r.Header.Get(CSRFHeaderName)

	// 2. Fall back to the POST (form) value.
	if issued == "" {
		issued = r.PostFormValue(CSRFFormFieldName)
		L(r).Debug("no CSRF token in header, reading POST form values",
			zap.String("content_type", r.Header.Get("Content-Type")))
	}

	// 3. Finally, fall back to the multipart form (if set).
	if issued == "" && r.MultipartForm != nil {
		L(r).Debug("no CSRF token in header or form, reading multi-part values")

		vals := r.MultipartForm.Value[CSRFFormFieldName]

		if len(vals) > 0 {
			issued = vals[0]
		}
	}

	// Decode the "issued" (pad + masked) token sent in the request. Return a
	// nil byte slice on a decoding error (this will fail upstream).
	decoded, err := base64.StdEncoding.DecodeString(issued)
	if err != nil {
		L(r).Debug("failed to base64 decode the CSRF token", zap.Error(err))
		return nil
	}

	return decoded
}

// Idempotent (safe) methods as defined by RFC7231 section 4.2.2.
var safeMethods = map[string]struct{}{
	http.MethodGet: {}, http.MethodHead: {}, http.MethodOptions: {}, http.MethodTrace: {}}

// gen will generate 'l' random bytes or panic.
func gen(l int) (b []byte) {
	b = make([]byte, l)
	if n, err := RandRead(b); n != l || err != nil {
		panic("failed to read CSRF random bytes: " + err.Error() + ", n:" + strconv.Itoa(n))
	}

	return
}

// mask returns a unique-per-request token to mitigate the BREACH attack
// as per http://breachattack.com/#mitigations
//
// The token is generated by XOR'ing a one-time-pad and the base (session) CSRF
// token and returning them together as a 64-byte slice. This effectively
// randomises the token on a per-request basis without breaking multiple browser
// tabs/windows.
func mask(realToken []byte) string {
	otp := gen(CSRFTokenLength)

	// XOR the OTP with the real token to generate a masked token. Append the
	// OTP to the front of the masked token to allow unmasking in the subsequent
	// request.
	return base64.StdEncoding.EncodeToString(append(otp, xorToken(otp, realToken)...))
}

// unmask splits the issued token (one-time-pad + masked token) and returns the
// unmasked request token for comparison.
func unmask(issued []byte) []byte {
	// Issued tokens are always masked and combined with the pad.
	if len(issued) != CSRFTokenLength*2 {
		return nil
	}

	// We now know the length of the byte slice.
	otp := issued[CSRFTokenLength:]
	masked := issued[:CSRFTokenLength]

	// Unmask the token by XOR'ing it against the OTP used to mask it.
	return xorToken(otp, masked)
}

// xorToken XORs tokens ([]byte) to provide unique-per-request CSRF tokens. It
// will return a masked token if the base token is XOR'ed with a one-time-pad.
// An unmasked token will be returned if a masked token is XOR'ed with the
// one-time-pad used to mask it.
func xorToken(a, b []byte) []byte {
	n := len(a)
	if len(b) < n {
		n = len(b)
	}

	res := make([]byte, n)

	for i := 0; i < n; i++ {
		res[i] = a[i] ^ b[i]
	}

	return res
}
