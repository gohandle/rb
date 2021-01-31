package rb

import (
	"errors"
	"reflect"

	"github.com/CloudyKit/jet/v6"
	"github.com/go-playground/validator/v10"
)

func (a *App) fieldErrorHelper(args jet.Arguments) (v reflect.Value) {
	args.RequireNumOfArguments("url", 2, 2)
	fname := args.Get(1).String()

	if (args.Get(0) == reflect.Value{}) {
		return reflect.ValueOf("") // passed in nil, also don't do anything
	}

	erri := args.Get(0).Interface()
	if erri == nil {
		return reflect.ValueOf("") // no error to display
	}

	err, ok := args.Get(0).Interface().(error)
	if !ok {
		args.Panicf("rb/field-error-helper: first argument must be an error type, got: %T", args.Get(0).Interface())
	}

	// try some commone interfaces that might give use a hint that the error is scoped to
	// a given field
	if ferr, ok := err.(interface{ Field() string }); ok && ferr.Field() == fname {
		return reflect.ValueOf(err)
	}

	// if it is a validator error, we can also determine the field
	var verr validator.ValidationErrors
	if errors.As(err, &verr) {
		for _, ferr := range verr {
			if ferr.StructField() == fname {
				return reflect.ValueOf(err)
			}
		}
	}

	// else, don't render
	return reflect.ValueOf("")
}
