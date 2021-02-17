package jethelper

import (
	"errors"
	"reflect"

	"github.com/CloudyKit/jet/v6"
	"github.com/go-playground/validator/v10"
)

// FieldError is a helper that will show an error if it implements a method
// that indicates that the error is for a certain field.
type FieldError jet.Func

// Name defines the name under which the helper is available in jet templates
func (FieldError) Name() string { return "rb_field_error" }

// NewFieldError actual helper
func NewFieldError() FieldError {
	return func(args jet.Arguments) (v reflect.Value) {
		args.RequireNumOfArguments(FieldError.Name(nil), 2, 2)
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

// NonFieldError is a helper that will show an error if it implements a method
// that indicates that the error is for a certain field.
type NonFieldError jet.Func

// Name defines the name under which the helper is available in jet templates
func (NonFieldError) Name() string { return "rb_non_field_error" }

// NewNonFieldError actual helper
func NewNonFieldError() NonFieldError {
	return func(args jet.Arguments) (v reflect.Value) {
		args.RequireNumOfArguments(NonFieldError.Name(nil), 1, 1)

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
