package rb

import (
	"fmt"
	"net/http"
)

type Bind interface {
	Bind(a *App, r *http.Request) error
	Value() (v interface{})
}

type ValueSource int

const (
	Form ValueSource = iota
	Query
	PostForm
)

type formBind struct {
	v interface{}
	s ValueSource
}

type FormBindOption func(*formBind)

func FromSource(s ValueSource) FormBindOption {
	return func(o *formBind) { o.s = s }
}

func FormBind(v interface{}, opts ...FormBindOption) Bind {
	b := formBind{v, Form}
	for _, opt := range opts {
		opt(&b)
	}

	return b
}

func (b formBind) Value() interface{} { return b.v }

func (b formBind) Bind(a *App, r *http.Request) error {
	if err := r.ParseForm(); err != nil {
		return err
	}

	switch b.s {
	case Form:
		return a.fdec.Decode(b.v, r.Form)
	case Query:
		return a.fdec.Decode(b.v, r.URL.Query())
	case PostForm:
		return a.fdec.Decode(b.v, r.PostForm)
	default:
		return fmt.Errorf("unsupported bind source: %T", b.s)
	}
}

type bindOptions struct {
	validate bool
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

	if err := b.Bind(a, r); err != nil {
		return fmt.Errorf("failed to bind: %v", err)
	}

	if !o.validate {
		return nil
	}

	if err := a.val.StructCtx(r.Context(), b.Value()); err != nil {
		return fmt.Errorf("failed to validate: %w", err)
	}

	return nil
}
