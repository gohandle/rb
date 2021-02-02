package rb

import (
	"net/http"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// HeaderRender is an optional interface that can be implemented by a Render. If it wants to
// change the header code it should return a new code and not call WriteHeader itself.
type HeaderRender interface {
	RenderHeader(a *App, w http.ResponseWriter, r *http.Request, status int) (int, error)
}

// Render can be implemented to customize how data should be written to a response.
type Render interface {
	Render(a *App, w http.ResponseWriter, r *http.Request) error
	Value() interface{}
}

type renderOpts struct {
	code int
}

func (r renderOpts) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddInt("status_code", r.code)
	return nil
}

// RenderOption configures a render
type RenderOption func(*renderOpts)

// Status configures any render such that the response header is written with the provided
// status code.
func Status(code int) RenderOption {
	return func(opts *renderOpts) {
		opts.code = code
	}
}

func (a *App) handleErrorOrPanic(w http.ResponseWriter, r *http.Request, err error) {
	a.L(r).Error("error while handling request", zap.Error(err))
	if a.ErrHandler == nil {
		a.L(r).Info("no error handler configured, render default error page")
		http.Error(w, "rb: no error handler but an error occured: "+err.Error(), http.StatusInternalServerError)
		return
	}

	herr := a.ErrHandler(a, w, r, err)
	if herr != nil {
		panic("rb: failed to handle error, original error: '" + err.Error() + "', error handler error: " + herr.Error())
	}
}

// Render can be called to write a representation of a value as the HTTP response. The application
// is responsible for creating a Render, for example usin JSON() for rendin json or Template()
// render a template.
//
// This method returns no error. Any error that occured during rendin will should be handled.
// It is possible to configure the ErrorHandler field on the app to configure application wide
// logic for handling errors durin rending.
func (a *App) Render(w http.ResponseWriter, r *http.Request, rr Render, opts ...RenderOption) {
	var o renderOpts
	for _, opt := range opts {
		opt(&o)
	}

	a.L(r).Debug("start render", zap.Any("render", rr), zap.Any("options", o))

	if hr, ok := rr.(HeaderRender); ok {
		var err error
		a.L(r).Debug("render implemented header rendering", zap.Int("status_code", o.code))
		o.code, err = hr.RenderHeader(a, w, r, o.code)
		if err != nil {
			a.L(r).Debug("error while rendering header", zap.Error(err))
			a.handleErrorOrPanic(w, r, err)
			return
		}
	}

	if o.code < 1 {
		o.code = http.StatusOK
		a.L(r).Debug("no explicit status provided, set default", zap.Int("status_code", o.code))
	}

	w.WriteHeader(o.code)
	a.L(r).Debug("wrote header", zap.Int("status_code", o.code))

	err := rr.Render(a, w, r)
	if err != nil {
		a.L(r).Debug("error while rendering body", zap.Error(err))
		a.handleErrorOrPanic(w, r, err)
		return
	}

	a.L(r).Debug("render complete")
}
