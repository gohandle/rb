package rb

import (
	"fmt"
	"net/http"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Bind interface can be implemented to support other data types for binding. This package implements
// form and JSON encoded request bodies.
type Bind interface {
	Bind(a *App, r *http.Request) error
	Value() (v interface{})
}

// ValueSource configures the source for a form decoding bind
type ValueSource int

const (
	// FormSource will decode from both the url query as the post body
	FormSource ValueSource = iota

	// QuerySource source will only bind the request's url query
	QuerySource

	// PostFormSource will only bind the request's form data in the body
	PostFormSource
)

type formBind struct {
	v interface{}
	s ValueSource
}

// FormBindOption configures the Form bind
type FormBindOption func(*formBind)

// FromSource configures from what part of the request the bind will decode form encoded values
func FromSource(s ValueSource) FormBindOption {
	return func(o *formBind) { o.s = s }
}

// Form will create a Bind that can be passed to the Bind method. This bind will decode request
// bodys that are application/www-x-form-url encoded and have that mime type as the Content-Type
// header.
func Form(v interface{}, opts ...FormBindOption) Bind {
	b := formBind{v, FormSource}
	for _, opt := range opts {
		opt(&b)
	}

	return b
}

func (b formBind) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddString("kind", "form")
	switch b.s {
	case FormSource:
		enc.AddString("source", "form")
	case QuerySource:
		enc.AddString("source", "query")
	case PostFormSource:
		enc.AddString("source", "post_form")
	}

	return nil
}

func (b formBind) Value() interface{} { return b.v }

func (b formBind) Bind(a *App, r *http.Request) error {
	if err := r.ParseForm(); err != nil {
		a.L(r).Error("failed to parse form", zap.Error(err))
		return err
	}

	switch b.s {
	case FormSource:
		a.L(r).Debug("binding from form", zap.Any("request_form", r.Form))
		return a.fdec.Decode(b.v, r.Form)
	case QuerySource:
		a.L(r).Debug("binding from query", zap.Any("request_urL_query", r.URL.Query()))
		return a.fdec.Decode(b.v, r.URL.Query())
	case PostFormSource:
		a.L(r).Debug("binding from post form", zap.Any("request_post_form", r.PostForm))
		return a.fdec.Decode(b.v, r.PostForm)
	default:
		return fmt.Errorf("unsupported bind source: %T", b.s)
	}
}

type bindOptions struct {
	validate bool
}

func (r bindOptions) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddBool("validate", r.validate)
	return nil
}

// BindOption configures any bind
type BindOption func(*bindOptions)

// AndValidate is a bind option that will cause the binding to also validate the bound value
// using the validator that is configured for this instance of *App.
func AndValidate() BindOption {
	return func(o *bindOptions) {
		o.validate = true
	}
}

// Bind will attemp to decode the body of request 'r' using bind 'b'. It takes an optional list
// of options to configure the binding process. It also supports validation of the bound value
// using the AndValidate() option.
func (a *App) Bind(r *http.Request, b Bind, opts ...BindOption) error {
	var o bindOptions
	for _, opt := range opts {
		opt(&o)
	}

	a.L(r).Debug("start bind",
		zap.Any("bind", b),
		zap.Any("options", o),
		zap.String("content_type", r.Header.Get("content-type")))

	if err := b.Bind(a, r); err != nil {
		a.L(r).Error("bind failed", zap.Error(err))
		return fmt.Errorf("failed to bind: %v", err)
	}

	if !o.validate {
		a.L(r).Debug("bind validate set to false, bind done")
		return nil
	}

	if err := a.val.StructCtx(r.Context(), b.Value()); err != nil {
		a.L(r).Error("validate failed", zap.Error(err))
		return fmt.Errorf("failed to validate: %w", err)
	}

	return nil
}
