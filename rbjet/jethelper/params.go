package jethelper

import (
	"reflect"

	"github.com/CloudyKit/jet/v6"
	rb "github.com/gohandle/rb"
)

// NewParams creates the jet template helper for retireving url parameters
func NewParams(rc rb.RouterCore) (string, jet.Func) {
	const name = "rb_url_param"
	return name, func(args jet.Arguments) (v reflect.Value) {
		args.RequireNumOfArguments(name, 1, 1)
		w, r, err := respReq(args)
		if err != nil {
			args.Panicf("failed to read response and request: %v", err)
		}

		return reflect.ValueOf(rc.Param(w, r, args.Get(0).String()))
	}
}
