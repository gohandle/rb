package rb

import "net/http"

type HeaderRender interface {
	RenderHeader(a *App, w http.ResponseWriter, r *http.Request, status int) error
}

type Render interface {
	Render(a *App, w http.ResponseWriter, r *http.Request) error
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

	if o.code < 1 {
		o.code = http.StatusOK
	}

	if hr, ok := rr.(HeaderRender); ok {
		if err := hr.RenderHeader(a, w, r, o.code); err != nil {
			panic("rb: failed to render header: " + err.Error())
		}
	} else {
		w.WriteHeader(o.code)
	}

	err := rr.Render(a, w, r)
	if err != nil {
		panic("rb: failed to render: " + err.Error())
	}
}
