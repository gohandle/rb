package jethelper

import (
	"reflect"

	"github.com/CloudyKit/jet/v6"
	"github.com/gohandle/rb/rbv2/rbcsrf"
)

// CSRF is a jet template helper that retrieves the current CSRF  token
type CSRF jet.Func

// Name defines the name under which the helper is available in jet templates
func (CSRF) Name() string { return "rb_csrf" }

// NewCSRF creates the jet template helper
func NewCSRF() CSRF {
	return func(args jet.Arguments) (v reflect.Value) {
		args.RequireNumOfArguments(CSRF.Name(nil), 0, 0)
		_, r, err := respReq(args)
		if err != nil {
			args.Panicf("failed to read response and request: %v", err)
		}

		tok := rbcsrf.Token(r.Context())
		if tok == "" {
			args.Panicf("no CSRF token in request context")
		}

		return reflect.ValueOf(tok)
	}
}
