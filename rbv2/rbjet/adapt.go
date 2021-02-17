package rbjet

import (
	"fmt"
	"net/http"
	"reflect"

	"github.com/CloudyKit/jet/v6"
	rb "github.com/gohandle/rb/rbv2"
)

const requestVarName = "rb_request"
const responseVarName = "rb_response"

type adapted struct{ *jet.Set }

// Adapt adapts a jset tmplate set to implement the template directory
func Adapt(jset *jet.Set, urlHelper URLHelper) rb.Templates {
	if urlHelper != nil {
		jset.AddGlobalFunc(urlHelper.Name(), jet.Func(urlHelper))
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

	vars[requestVarName] = reflect.ValueOf(r)
	vars[responseVarName] = reflect.ValueOf(w)

	return e.Template.Execute(w, vars, data)
}

// respReq is a utility method that get the request and response from variables in the execution context
func respReq(a jet.Arguments) (w http.ResponseWriter, r *http.Request, err error) {
	reqv := a.Runtime().Resolve(requestVarName)
	if (reqv == reflect.Value{}) {
		return nil, nil, fmt.Errorf("failed to resolve '%s' variable", requestVarName)
	}

	r, ok := reqv.Interface().(*http.Request)
	if !ok {
		return nil, nil, fmt.Errorf("failed to turn '%s' variable to *http.Request", requestVarName)
	}

	respv := a.Runtime().Resolve(responseVarName)
	if (reqv == reflect.Value{}) {
		return nil, nil, fmt.Errorf("failed to resolve '%s' variable", responseVarName)
	}

	w, ok = respv.Interface().(http.ResponseWriter)
	if !ok {
		return nil, nil, fmt.Errorf("failed to turn '%s' variable to http.ResponseWriter", responseVarName)
	}

	return
}
