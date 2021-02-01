package rb

import (
	"net/http"
)

// HeaderRender is an optional interface that can be implemented by a Render. If it wants to
// change the header code it should return a new code and not call WriteHeader itself.
type HeaderRender interface {
	RenderHeader(a *App, w http.ResponseWriter, r *http.Request, status int) (int, error)
}

type Render interface {
	Render(a *App, w http.ResponseWriter, r *http.Request) error
	Value() interface{}
}

type renderOpts struct {
	code int
}

type RenderOption func(*renderOpts)

func Status(code int) RenderOption {
	return func(opts *renderOpts) {
		opts.code = code
	}
}

func (a *App) handleErrorOrPanic(w http.ResponseWriter, r *http.Request, err error) {
	// @TODO log error
	if a.ErrHandler == nil {
		http.Error(w, "rb: no error handler but an error occured: "+err.Error(), http.StatusInternalServerError)
		return
	}

	herr := a.ErrHandler(a, w, r, err)
	if herr != nil {
		panic("rb: failed to handle error, original error: '" + err.Error() + "', error handler error: " + herr.Error())
	}
	// @TODO log error handling
}

func (a *App) Render(w http.ResponseWriter, r *http.Request, rr Render, opts ...RenderOption) {
	var o renderOpts
	for _, opt := range opts {
		opt(&o)
	}

	if hr, ok := rr.(HeaderRender); ok {
		var err error
		o.code, err = hr.RenderHeader(a, w, r, o.code)
		if err != nil {
			a.handleErrorOrPanic(w, r, err)
			return
		}
	}

	if o.code < 1 {
		o.code = http.StatusOK
	}

	w.WriteHeader(o.code)
	err := rr.Render(a, w, r)
	if err != nil {
		a.handleErrorOrPanic(w, r, err)
		return
	}
}
