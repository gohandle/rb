package rb

import (
	"net/http"

	"go.uber.org/zap"
)

// ActionFunc implements an application action
type ActionFunc func(Ctx) error

// App provides application wide functionality
type App interface {
	Action(af ActionFunc) http.Handler
}

// DefaultApp contains sensible dependencies to
// create a App quickly
type DefaultApp struct {
	core Core
}

// New creates an app using default dependencies
func New(core Core) App {
	return &DefaultApp{core}
}

// Action creates an http.Handler from our action func
func (a *DefaultApp) Action(af ActionFunc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c := NewCtx(w, r, a.core)
		if err := af(c); err != nil {
			a.core.HandleError(w, c.Request(), err)
		}
	})
}

// L is utility method that returns a zap logger, if possible from the request
// context. If not it will return a NopLogger
func L(r ...*http.Request) *zap.Logger {
	if len(r) < 1 {
		return zap.NewNop()
	}

	l := RequestLogger(r[0].Context())
	if l != nil {
		return l
	}

	l, _ = zap.NewDevelopment()
	return l
}
