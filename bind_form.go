package rb

import (
	"fmt"
	"net/http"

	"go.uber.org/zap"
)

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

// FromSource configures from what part of the request the bind will decode form encoded values
func FromSource(s ValueSource) FormBindOption {
	return func(o *formBind) { o.s = s }
}

// FormBindOption configures the Form bind
type FormBindOption func(*formBind)

type formBind struct {
	v interface{}
	s ValueSource
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

func (b formBind) Value() interface{} { return b.v }

func (b formBind) Bind(bc BindCore, r *http.Request) error {

	if err := r.ParseForm(); err != nil {
		L(r).Error("failed to parse form", zap.Error(err))
		return err
	}

	switch b.s {
	case FormSource:
		L(r).Debug("binding from form", zap.Any("request_form", r.Form))
		return bc.DecodeValues(b.v, r.Form)
	case QuerySource:
		L(r).Debug("binding from query", zap.Any("request_urL_query", r.URL.Query()))
		return bc.DecodeValues(b.v, r.URL.Query())
	case PostFormSource:
		L(r).Debug("binding from post form", zap.Any("request_post_form", r.PostForm))
		return bc.DecodeValues(b.v, r.PostForm)
	default:
		return fmt.Errorf("unsupported bind source: %T", b.s)
	}
}
