package rb

import "net/http"

type Render interface {
	Execute(w http.ResponseWriter, r *http.Request) error
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

func (a *App) Render(w http.ResponseWriter, r *http.Request, rr Render, opts ...RenderOption) {
	var o renderOpts
	for _, opt := range opts {
		opt(&o)
	}

	err := rr.Execute(w, r)
	if err != nil {
		// @TODO handle error
	}
}
