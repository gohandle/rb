package rb

import (
	"fmt"
	"net/http"
	"reflect"

	"github.com/CloudyKit/jet/v6"
)

// TemplateRenderOption allows for configurin the rendering of a template
type TemplateRenderOption func(*templateRender)

// TemplateVar configures the template render to make an extra variable available in the scope of
// the template.
func TemplateVar(name string, v interface{}) TemplateRenderOption {
	return func(r *templateRender) {
		if r.vars == nil {
			r.vars = make(jet.VarMap)
		}

		r.vars[name] = reflect.ValueOf(v)
	}
}

type templateRender struct {
	name string
	val  interface{}
	vars map[string]reflect.Value
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

func (r templateRender) Render(rc RenderCore, wr http.ResponseWriter, req *http.Request) error {
	tmpl, err := rc.Lookup(r.name)
	if err != nil {
		return fmt.Errorf("failed to get template: %w", err)
	}

	return tmpl.Execute(wr, req, r.vars, r.val)
}
