package rb

import (
	"fmt"
	"net/http"

	"github.com/CloudyKit/jet/v6"
	"go.uber.org/zap/zapcore"
)

// TemplateRenderOption allos for configurin the rendering of a template
type TemplateRenderOption func(*templateRender)

type templateRender struct {
	name string
	val  interface{}
	vars jet.VarMap
}

func (r templateRender) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddString("kind", "template")
	enc.AddString("template", r.name)
	for k, v := range r.vars {
		enc.AddString("var_"+k, v.Kind().String())
	}

	return nil
}

func (r templateRender) Value() interface{} { return r.val }

func (r templateRender) Render(a *App, wr http.ResponseWriter, req *http.Request) error {
	tmpl, err := a.view.GetTemplate(r.name)
	if err != nil {
		return fmt.Errorf("failed to get template: %w", err)
	}

	return tmpl.Execute(wr, r.vars, r.val)
}

// TemplateVar configures the template render to make an extra variable available in the scope of
// the template.
func TemplateVar(name string, v interface{}) TemplateRenderOption {
	return func(r *templateRender) {
		if r.vars == nil {
			r.vars = make(jet.VarMap)
		}

		r.vars.Set(name, v)
	}
}

// Template will create a render that uses the jet templating engine to render HTML. It will
// take 'v' as the value that is rendered and makes it availalbe as '.' in the template.
func Template(name string, data interface{}, opts ...TemplateRenderOption) Render {
	r := templateRender{
		name: name, val: data}
	for _, opt := range opts {
		opt(&r)
	}
	return r
}
