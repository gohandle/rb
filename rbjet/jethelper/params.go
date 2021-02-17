package jethelper

import (
	"reflect"

	"github.com/CloudyKit/jet/v6"
	rb "github.com/gohandle/rb"
)

// Params is a jet template helper that retrieves url parameters
type Params jet.Func

// Name defines the name under which the helper is available in jet templates
func (Params) Name() string { return "rb_params" }

// NewParams creates the jet template helper for retireving url parameters
func NewParams(rc rb.RouterCore) Params {
	return func(args jet.Arguments) (v reflect.Value) {
		args.RequireNumOfArguments(Params.Name(nil), 0, 0)
		w, r, err := respReq(args)
		if err != nil {
			args.Panicf("failed to read response and request: %v", err)
		}

		return reflect.ValueOf(rc.Params(w, r))
	}
}
