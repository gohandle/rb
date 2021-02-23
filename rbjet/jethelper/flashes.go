package jethelper

import (
	"reflect"

	"github.com/CloudyKit/jet/v6"
	rb "github.com/gohandle/rb"
)

// NewFlashes creates the jet template helper
func NewFlashes() (string, jet.Func) {
	const name = "rb_flashes"
	return name, func(args jet.Arguments) (v reflect.Value) {
		args.RequireNumOfArguments(name, 0, 0)
		_, r, err := respReq(args)
		if err != nil {
			args.Panicf("failed to read response and request: %v", err)
		}

		return reflect.ValueOf(rb.FlashMessages(r.Context()))
	}
}
