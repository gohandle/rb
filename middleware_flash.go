package rb

import (
	"context"
	"net/http"

	"go.uber.org/zap"
)

// FlashMiddleware pops flash messages from the session and stores it in the context
type FlashMiddleware func(http.Handler) http.Handler

// FlashSessionField is the field name of the sesion that contains flash messages
const FlashSessionField = "_rb_flash"

// FlashMessages returns the id token from the context if it has any
func FlashMessages(ctx context.Context) (msgs []string) {
	msgs, _ = ctx.Value(ctxKey("flash")).([]string)
	return
}

// WithFlashMessages sets the id token context value
func WithFlashMessages(ctx context.Context, msgs []string) context.Context {
	return context.WithValue(ctx, ctxKey("flash"), msgs)
}

// Flash is a small helper that appends a flash message to the session. The Ctx
// implements the first argument.
func Flash(c interface {
	Session(o ...SessionOption) Session
}, ms ...string) {
	msgs, _ := c.Session().Get(FlashSessionField).([]string)
	msgs = append(msgs, ms...)
	c.Session().Set(FlashSessionField, msgs)
}

// Flashes returns any flash messages that were present at the beginning of the
// request. It does not include any flash messages that were set during the handling
// of the request. It is a shortcut for calling FlashMessages with the request's context.
func Flashes(c interface {
	Request() *http.Request
}) []string {
	return FlashMessages(c.Request().Context())
}

// NewFlashMiddleware creates the actual middleware
func NewFlashMiddleware(sc SessionCore) LoggerMiddleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			v := sc.Session(w, r).Pop(FlashSessionField)
			msgs, ok := v.([]string)
			L(r).Debug("called pop on session for flash messages",
				zap.String("session_field", FlashSessionField),
				zap.Bool("field_exists", v != nil),
				zap.Bool("correct_type", ok), zap.Int("num_msgs", len(msgs)))

			if ok && len(msgs) > 0 {

				r = r.WithContext(WithFlashMessages(r.Context(), msgs))
			}

			next.ServeHTTP(w, r)
		})
	}
}
