package jethelper

import (
	"reflect"

	"github.com/CloudyKit/jet/v6"
	rb "github.com/gohandle/rb"
)

// Flashes is a jet template helper that retrieves any flashes that were read at the beginning
// of the request handling
type Flashes jet.Func

// Name defines the name under which the helper is available in jet templates
func (Flashes) Name() string { return "rb_flashes" }

// NewFlashes creates the jet template helper
func NewFlashes() Flashes {
	return func(args jet.Arguments) (v reflect.Value) {
		args.RequireNumOfArguments(Flashes.Name(nil), 0, 0)
		_, r, err := respReq(args)
		if err != nil {
			args.Panicf("failed to read response and request: %v", err)
		}

		return reflect.ValueOf(rb.FlashMessages(r.Context()))
	}
}
