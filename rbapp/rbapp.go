package rbapp

import (
	"net/http"

	"github.com/gohandle/rb"
	"github.com/gohandle/rb/rbjet/jethelper"
	"github.com/gohandle/rb/rbjit"
	"go.uber.org/zap"
)

// DefaultApp contains sensible default dependencies to create a App quickly
type DefaultApp struct {
	core rb.Core
}

// New creates an app using the provided core while adding the default Middleware to the router core.
func New(core rb.Core, logs *zap.Logger) rb.App {
	core.Use(rb.NewIDMiddleware(rb.CommonRequestIDHeaders...))
	core.Use(rb.NewLoggerMiddleware(logs))
	core.Use(rbjit.NewMiddleware())
	core.Use(rb.NewSessionSaveMiddleware(core))
	core.Use(rb.NewCSRFMiddlware(core, core))
	core.Use(rb.NewFlashMiddleware(core))

	core.AddHelper(jethelper.NewCSRF())
	core.AddHelper(jethelper.NewFieldError())
	core.AddHelper(jethelper.NewNonFieldError())
	core.AddHelper(jethelper.NewFlashes())
	core.AddHelper(jethelper.NewParams(core))
	core.AddHelper(jethelper.NewRoute(core))
	core.AddHelper(jethelper.NewSession(core))
	core.AddHelper(jethelper.NewTrans(core))
	core.AddHelper(jethelper.NewURL(core))

	return &DefaultApp{core}
}

// Action creates an http.Handler from our action func
func (a *DefaultApp) Action(af rb.ActionFunc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c := rb.NewCtx(rbjit.New(w), r, a.core)
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

	l := rb.RequestLogger(r[0].Context())
	if l != nil {
		return l
	}

	l, _ = zap.NewDevelopment()
	return l
}
