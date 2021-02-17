package jethelper

import (
	"reflect"

	"github.com/CloudyKit/jet/v6"
	rb "github.com/gohandle/rb/rbv2"
)

// Session is a jet template helper that retrieves url parameters
type Session jet.Func

// Name defines the name under which the helper is available in jet templates
func (Session) Name() string { return "rb_session" }

// NewSession creates the jet template helper for retireving url parameters
func NewSession(sc rb.SessionCore) Session {
	return func(args jet.Arguments) (v reflect.Value) {
		args.RequireNumOfArguments(Session.Name(nil), 1, 1)
		w, r, err := respReq(args)
		if err != nil {
			args.Panicf("failed to read response and request: %v", err)
		}

		k := args.Get(0).Interface()
		if k == nil {
			args.Panicf("requires a non-nil value as the first argument, got: '%v'", k)
		}

		return reflect.ValueOf(sc.Session(w, r).Get(k))
	}
}
