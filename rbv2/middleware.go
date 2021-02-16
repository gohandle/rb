package rb

import (
	"context"
	"encoding/base64"
	"math/rand"
	"net/http"

	"go.uber.org/zap"
)

// make context keys private
type ctxKey string

// RequestID returns the id token from the context if it has any
func RequestID(ctx context.Context) (s string) {
	s, _ = ctx.Value(ctxKey("id")).(string)
	return
}

// WithRequestID sets the id token context value
func WithRequestID(ctx context.Context, s string) context.Context {
	return context.WithValue(ctx, ctxKey("id"), s)
}

// CommonRequestIDHeaders can be provided to the middleware to support some common
// headers for request identification.
var CommonRequestIDHeaders = []string{
	"X-Request-ID", "X-Correlation-ID", // unofficial standards: https://en.wikipedia.org/wiki/List_of_HTTP_header_fields
	"X-Amzn-Trace-Id",       // AWS xray tracing: https://docs.aws.amazon.com/xray/latest/devguide/xray-concepts.html#xray-concepts-tracingheader
	"Cf-Request-Id",         // Cloudflare: https://community.cloudflare.com/t/new-http-response-header-cf-request-id/165869
	"X-Cloud-Trace-Context", // Google Cloud https://cloud.google.com/appengine/docs/standard/go/reference/request-response-headers
}

// RandRead is used for request id generation. It can be ovewritten in test to make them fully
// deterministic. The default is set to a non-cryptographic random number
var RandRead = rand.Read

// IDMiddleware creates middleware that looks at common request identification headers and
// makes it available to the request's context. If none of the request headers are provided
// a new id is generated.
func IDMiddleware(hdrs ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var rid string
			for _, hdrn := range hdrs {
				rid = r.Header.Get(hdrn)
				if rid != "" {
					break
				}
			}

			if rid == "" {
				var b [18]byte
				if _, err := RandRead(b[:]); err != nil {
					L(r).Error("failed to read random bytes for request id middleware",
						zap.Error(err))
				}
				rid = base64.URLEncoding.EncodeToString(b[:])
			}

			next.ServeHTTP(w, r.WithContext(
				WithRequestID(r.Context(), rid)))
		})
	}
}

// WithRequestLogger sets request scoped logger
func WithRequestLogger(ctx context.Context, logs *zap.Logger) context.Context {
	return context.WithValue(ctx, ctxKey("logger"), logs)
}

// RequestLogger returns the request scoped logger, returns nil if none is configured
func RequestLogger(ctx context.Context) (l *zap.Logger) {
	l, _ = ctx.Value(ctxKey("logger")).(*zap.Logger)
	return
}

// LoggerMiddleware will create a request scoped logger that uses the request id to make those logs
// observable for debugging.
func LoggerMiddleware(logs *zap.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			l := logs.With(
				zap.String("request_url", r.URL.String()),
				zap.String("request_method", r.Method),
				zap.Int64("request_length", r.ContentLength),
				zap.String("request_host", r.Host),
				zap.String("request_uri", r.RequestURI),
			)

			if rid := RequestID(r.Context()); rid != "" {
				l = l.With(zap.String("request_id", rid))
			}

			next.ServeHTTP(w, r.WithContext(
				WithRequestLogger(r.Context(), l)))
		})
	}
}
