package rb

import "net/http"

type InjectorFunc func(a *App, v interface{}) error

func (f InjectorFunc) OnRender(a *App, v interface{}) error {
	return f(a, v)
}

type Injector interface {
	OnRender(a *App, v interface{}) error
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
	if err := r.inj.OnRender(a, r.val); err != nil {
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
