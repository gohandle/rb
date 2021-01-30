package rb

import (
	"fmt"
	"net/http"

	"github.com/CloudyKit/jet/v6"
)

type TemplateRenderOption func(*templateRender)

type templateRender struct {
	name string
	data interface{}
	vars jet.VarMap
}

func (r templateRender) Render(a *App, wr http.ResponseWriter, req *http.Request) error {
	tmpl, err := a.view.GetTemplate(r.name)
	if err != nil {
		return fmt.Errorf("failed to get template: %w", err)
	}

	return tmpl.Execute(wr, r.vars, r.data)
}

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
