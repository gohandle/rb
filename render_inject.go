package rb

import "net/http"

type InjectorFunc func(a *App, w http.ResponseWriter, req *http.Request, v interface{}) error

func (f InjectorFunc) OnRender(a *App, w http.ResponseWriter, req *http.Request, v interface{}) error {
	return f(a, w, req, v)
}

// Injector can be implemented to allow hooks to be called just before the body is rendered. The
// ResponseWriter is provided to allow for the reading of headers, but at this point in the
// lifecycle the header is already written and cannot be changed.
type Injector interface {
	OnRender(a *App, w http.ResponseWriter, req *http.Request, v interface{}) error
}

type renderInject struct {
	inj Injector
	val interface{}
	rr  Render
	hr  HeaderRender
}

func (r renderInject) RenderHeader(a *App, w http.ResponseWriter, req *http.Request, status int) (int, error) {
	if r.hr == nil {
		return status, nil // no-op header render
	}

	return r.hr.RenderHeader(a, w, req, status)
}

func (r renderInject) Render(a *App, wr http.ResponseWriter, req *http.Request) error {
	if err := r.inj.OnRender(a, wr, req, r.val); err != nil {
		return err // injector failed
	}

	return r.rr.Render(a, wr, req)
}

func (r renderInject) Value() interface{} { return r.val }

func Inject(rr Render, inj Injector) Render {
	r := renderInject{rr: rr, val: rr.Value(), inj: inj}
	if hr, ok := rr.(HeaderRender); ok {
		r.hr = hr
	}

	return r
}
