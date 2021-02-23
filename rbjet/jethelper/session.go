package jethelper

import (
	"reflect"

	"github.com/CloudyKit/jet/v6"
	rb "github.com/gohandle/rb"
)

// NewSession creates the jet template helper for reading session data
func NewSession(sc rb.SessionCore) (string, jet.Func) {
	const name = "rb_session"
	return name, func(args jet.Arguments) (v reflect.Value) {
		args.RequireNumOfArguments(name, 1, 1)
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
