package rbjet

import (
	"net/http"
	"reflect"

	"github.com/CloudyKit/jet/v6"
	rb "github.com/gohandle/rb/rbv2"
	"github.com/gohandle/rb/rbv2/rbjet/jethelper"
)

type adapted struct{ *jet.Set }

// Adapt adapts a jset tmplate set to implement the template directory
func Adapt(
	jset *jet.Set,
	urlHelper jethelper.URL,
	paramsHelper jethelper.Params,
	routeHelper jethelper.Route,
	transHelper jethelper.Trans,
	sessionHelper jethelper.Session,
) rb.Templates {
	if urlHelper != nil {
		jset.AddGlobalFunc(urlHelper.Name(), jet.Func(urlHelper))
	}

	if paramsHelper != nil {
		jset.AddGlobalFunc(paramsHelper.Name(), jet.Func(paramsHelper))
	}

	if routeHelper != nil {
		jset.AddGlobalFunc(routeHelper.Name(), jet.Func(routeHelper))
	}

	if transHelper != nil {
		jset.AddGlobalFunc(transHelper.Name(), jet.Func(transHelper))
	}

	if sessionHelper != nil {
		jset.AddGlobalFunc(sessionHelper.Name(), jet.Func(sessionHelper))
	}

	return adapted{jset}
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
