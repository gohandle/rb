package rbjet

import (
	"net/http"
	"reflect"

	"github.com/CloudyKit/jet/v6"
	rb "github.com/gohandle/rb"
	"github.com/gohandle/rb/rbjet/jethelper"
)

type adapted struct{ *jet.Set }

// Adapt adapts a jset tmplate set to implement the template directory
func Adapt(
	jset *jet.Set,
) rb.Templates {
	return adapted{jset}
}

func (a adapted) AddHelper(name string, v interface{}) rb.Templates {
	hf, ok := v.(jet.Func)
	if !ok {
		panic("rbjet: called add helper with invalid value. Must be a jet.Func")
	}

	a.Set.AddGlobalFunc(name, hf)
	return a
}

func (a adapted) Lookup(n string) (rb.TemplateExecuter, error) {
	jset, err := a.Set.GetTemplate(n)
	if err != nil {
		return nil, err
	}

	return exec{jset}, nil
}

type exec struct{ *jet.Template }

func (e exec) Execute(
	w http.ResponseWriter,
	r *http.Request,
	vars map[string]reflect.Value,
	data interface{},
) (err error) {
	if vars == nil {
		vars = make(map[string]reflect.Value)
	}

	vars[jethelper.RequestVarName] = reflect.ValueOf(r)
	vars[jethelper.ResponseVarName] = reflect.ValueOf(w)

	return e.Template.Execute(w, vars, data)
}
