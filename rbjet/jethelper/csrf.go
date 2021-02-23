package jethelper

import (
	"reflect"

	"github.com/CloudyKit/jet/v6"
	"github.com/gohandle/rb"
)

// NewCSRF creates the jet template helper
func NewCSRF() (string, jet.Func) {
	const name = "rb_csrf"
	return name, func(args jet.Arguments) (v reflect.Value) {
		args.RequireNumOfArguments(name, 0, 0)
		_, r, err := respReq(args)
		if err != nil {
			args.Panicf("failed to read response and request: %v", err)
		}

		tok := rb.CSRFToken(r.Context())
		if tok == "" {
			args.Panicf("no CSRF token in request context")
		}

		return reflect.ValueOf(tok)
	}
}
