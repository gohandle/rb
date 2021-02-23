package jethelper

import (
	"reflect"

	"github.com/CloudyKit/jet/v6"
	rb "github.com/gohandle/rb"
)

// NewURL creates the jet template helper for generating urls. It requires that the
// template has variables for the html request, and response.
func NewURL(rc rb.RouterCore) (string, jet.Func) {
	const name = "rb_url"
	return name, func(args jet.Arguments) (v reflect.Value) {
		args.RequireNumOfArguments(name, 1, -1)

		w, r, err := respReq(args)
		if err != nil {
			args.Panicf("failed to read response and request: %v", err)
		}

		s := args.Get(0).String()
		if s == "" {
			args.Panicf("requires a non-empty string as the first argument, got: '%v'", s)
		}

		var opts []rb.URLOption
		for i := 2; i < args.NumOfArguments(); i += 2 {
			if args.Get(i-1).Kind() != reflect.String || args.Get(i).Kind() != reflect.String {
				args.Panicf("parameter arguments bust be strings, got: k:%v/v:%v", args.Get(i-1).Kind(), args.Get(i).Kind())
			}

			k, v := args.Get(i-1).String(), args.Get(i).String()
			opts = append(opts, rb.URLVar(k, v))
		}

		return reflect.ValueOf(rc.URL(w, r, s, opts...))
	}
}
