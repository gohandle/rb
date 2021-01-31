package rb

import (
	"reflect"

	"github.com/CloudyKit/jet/v6"
)

func (a *App) urlHelper(args jet.Arguments) (v reflect.Value) {
	args.RequireNumOfArguments("url", 1, -1)
	s := args.Get(0).String()
	if s == "" {
		args.Panicf("rb/url-helper: requires a non-empty string as the first argument, got: '%v'", s)
	}

	var opts []URLOption
	for i := 2; i < args.NumOfArguments(); i += 2 {
		if args.Get(i-1).Kind() != reflect.String || args.Get(i).Kind() != reflect.String {
			args.Panicf("rb/url-helper: parameter arguments bust be strings, got: k:%v/v:%v", args.Get(i-1).Kind(), args.Get(i).Kind())
		}

		k, v := args.Get(i-1).String(), args.Get(i).String()
		opts = append(opts, URLVar(k, v))
	}

	loc, err := a.GenerateURL(s, opts...)
	if err != nil {
		args.Panicf("rb/url-helper: failed to generate url: %v", err)
	}

	return reflect.ValueOf(loc)
}
