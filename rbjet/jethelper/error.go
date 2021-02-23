package jethelper

import (
	"errors"
	"reflect"

	"github.com/CloudyKit/jet/v6"
	"github.com/go-playground/validator/v10"
)

// NewFieldError actual helper
func NewFieldError() (string, jet.Func) {
	const name = "rb_field_error"
	return name, func(args jet.Arguments) (v reflect.Value) {
		args.RequireNumOfArguments(name, 2, 2)
		fname := args.Get(1).String()

		ok, err := castErrArg(args.Get(0))
		if !ok {
			args.Panicf("first argument must be an error type, got: %T", args.Get(0).Interface())
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

		return
	}
}

// NewNonFieldError actual helper
func NewNonFieldError() (string, jet.Func) {
	const name = "rb_non_field_error"
	return name, func(args jet.Arguments) (v reflect.Value) {
		args.RequireNumOfArguments(name, 1, 1)

		ok, err := castErrArg(args.Get(0))
		if !ok {
			args.Panicf("first argument must be an error type, got: %T", args.Get(0).Interface())
		}

		if _, ok := err.(interface{ Field() string }); ok {
			return v
		}

		var verr validator.ValidationErrors
		if errors.As(err, &verr) {
			return v
		}

		return reflect.ValueOf(err)
	}
}

// castErrArg will try to cas a jet argument into an erro
func castErrArg(rv reflect.Value) (bool, error) {
	if (rv == reflect.Value{}) || rv.Interface() == nil {
		return true, nil
	}

	err, ok := rv.Interface().(error)
	if !ok {
		return false, nil
	}

	return true, err
}
