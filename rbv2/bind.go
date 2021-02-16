package rb

import (
	"fmt"
	"net/http"

	"go.uber.org/zap"
)

// Bind interface can be implemented to support other data types for binding. This package implements
// form and JSON encoded request bodies.
type Bind interface {
	Bind(bc BindCore, r *http.Request) error
	Value() (v interface{})
}

// BindOptions hold the resulting binding options
type BindOptions struct {
	validate bool
}

// BindOption configures any bind
type BindOption func(*BindOptions)

// AndValidate is a bind option that will cause the binding to also validate the bound value
// using the validator that is configured for this instance of *App.
func AndValidate() BindOption {
	return func(o *BindOptions) {
		o.validate = true
	}
}

type bindCore struct {
	ValuesDecoder
	StructValidator
}

// NewBindCore inits the bind core part of the core
func NewBindCore(fdec ValuesDecoder, val StructValidator) BindCore {
	return &bindCore{
		fdec,
		val,
	}
}

func (bc *bindCore) Bind(w http.ResponseWriter, r *http.Request, b Bind, opts ...BindOption) (bool, error) {
	var o BindOptions
	for _, opt := range opts {
		opt(&o)
	}

	L(r).Debug("start bind",
		zap.Any("bind", b),
		zap.Any("options", o),
		zap.String("content_type", r.Header.Get("content-type")))

	if err := b.Bind(bc, r); err != nil {
		L(r).Error("bind failed", zap.Error(err))
		return true, fmt.Errorf("failed to bind: %v", err)
	}

	if !o.validate {
		L(r).Debug("bind validate set to false, bind done")
		return true, nil
	}

	if err := bc.ValidateStruct(r.Context(), b.Value()); err != nil {
		L(r).Error("validate failed", zap.Error(err))
		return true, fmt.Errorf("failed to validate: %w", err)
	}

	return true, nil
}
