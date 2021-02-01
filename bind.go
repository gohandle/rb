package rb

import (
	"fmt"
	"net/http"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Bind interface {
	Bind(a *App, r *http.Request) error
	Value() (v interface{})
}

type ValueSource int

const (
	FormSource ValueSource = iota
	QuerySource
	PostFormSource
)

type formBind struct {
	v interface{}
	s ValueSource
}

type FormBindOption func(*formBind)

func FromSource(s ValueSource) FormBindOption {
	return func(o *formBind) { o.s = s }
}

func Form(v interface{}, opts ...FormBindOption) Bind {
	b := formBind{v, FormSource}
	for _, opt := range opts {
		opt(&b)
	}

	return b
}

func (r formBind) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddString("kind", "form")
	switch r.s {
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
		a.L(r).Error("binding from form", zap.Any("request_form", r.Form))
		return a.fdec.Decode(b.v, r.Form)
	case QuerySource:
		a.L(r).Error("binding from query", zap.Any("request_urL_query", r.URL.Query()))
		return a.fdec.Decode(b.v, r.URL.Query())
	case PostFormSource:
		a.L(r).Error("binding from post form", zap.Any("request_post_form", r.PostForm))
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

type BindOption func(*bindOptions)

func AndValidate() BindOption {
	return func(o *bindOptions) {
		o.validate = true
	}
}

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
