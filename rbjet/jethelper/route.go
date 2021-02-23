package jethelper

import (
	"reflect"

	"github.com/CloudyKit/jet/v6"
	rb "github.com/gohandle/rb"
)

// NewRoute creates the jet template helper for retrieves the current route
func NewRoute(rc rb.RouterCore) (string, jet.Func) {
	const name = "rb_route"
	return name, func(args jet.Arguments) (v reflect.Value) {
		args.RequireNumOfArguments(name, 0, 0)
		w, r, err := respReq(args)
		if err != nil {
			args.Panicf("failed to read response and request: %v", err)
		}

		return reflect.ValueOf(rc.Route(w, r))
	}
}
