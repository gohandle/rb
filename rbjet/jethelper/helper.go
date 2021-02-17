package jethelper

import (
	"fmt"
	"net/http"
	"reflect"

	"github.com/CloudyKit/jet/v6"
)

// RequestVarName is the variable holding the http request
const RequestVarName = "rb_request"

// ResponseVarName is the variable holding the http response writer
const ResponseVarName = "rb_response"

// respReq is a utility method that get the request and response from variables in the execution context
func respReq(a jet.Arguments) (w http.ResponseWriter, r *http.Request, err error) {
	reqv := a.Runtime().Resolve(RequestVarName)
	if (reqv == reflect.Value{}) {
		return nil, nil, fmt.Errorf("failed to resolve '%s' variable", RequestVarName)
	}

	r, ok := reqv.Interface().(*http.Request)
	if !ok {
		return nil, nil, fmt.Errorf("failed to turn '%s' variable to *http.Request", RequestVarName)
	}

	respv := a.Runtime().Resolve(ResponseVarName)
	if (reqv == reflect.Value{}) {
		return nil, nil, fmt.Errorf("failed to resolve '%s' variable", ResponseVarName)
	}

	w, ok = respv.Interface().(http.ResponseWriter)
	if !ok {
		return nil, nil, fmt.Errorf("failed to turn '%s' variable to http.ResponseWriter", ResponseVarName)
	}

	return
}
