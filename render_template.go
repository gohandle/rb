package rb

import (
	"net/http"

	"github.com/CloudyKit/jet/v6"
)

type TemplateRenderOption func(*templateRender)

type templateRender struct {
	name string
	data interface{}
	vars jet.VarMap
}

func (r templateRender) Execute(wr http.ResponseWriter, req *http.Request) error { return nil }

func TemplateVar(name string, v interface{}) TemplateRenderOption {
	return func(r *templateRender) {
		if r.vars == nil {
			r.vars = make(jet.VarMap)
		}

		r.vars.Set(name, v)
	}
}

func Template(name string, data interface{}, opts ...TemplateRenderOption) Render {
	r := templateRender{
		name: name, data: data}
	for _, opt := range opts {
		opt(&r)
	}
	return r
}
