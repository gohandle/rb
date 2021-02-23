package jethelper

import (
	"reflect"

	"github.com/CloudyKit/jet/v6"
	rb "github.com/gohandle/rb"
)

// NewTrans creates the jet template helper for retireving url parameters
func NewTrans(tc rb.TranslateCore) (string, jet.Func) {
	const name = "t"
	return name, func(args jet.Arguments) (v reflect.Value) {
		args.RequireNumOfArguments(name, 1, 2)
		w, r, err := respReq(args)
		if err != nil {
			args.Panicf("failed to read response and request: %v", err)
		}

		mid := args.Get(0).String()
		if mid == "" {
			args.Panicf("requires a non-empty string as the first argument, got: '%v'", mid)
		}

		return reflect.ValueOf(tc.Translate(w, r, mid))
	}
}
