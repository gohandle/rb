package jethelper

import (
	"reflect"

	"github.com/CloudyKit/jet/v6"
	rb "github.com/gohandle/rb/rbv2"
)

// Route is a jet template helper that retrieves the current route
type Route jet.Func

// Name defines the name under which the helper is available in jet templates
func (Route) Name() string { return "rb_route" }

// NewRoute creates the jet template helper for retrieves the current route
func NewRoute(rc rb.RouterCore) Route {
	return func(args jet.Arguments) (v reflect.Value) {
		args.RequireNumOfArguments(Route.Name(nil), 0, 0)
		w, r, err := respReq(args)
		if err != nil {
			args.Panicf("failed to read response and request: %v", err)
		}

		return reflect.ValueOf(rc.Route(w, r))
	}
}
