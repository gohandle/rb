package rb

import (
	"context"
	"io"
	"net/http"
	"net/url"
	"reflect"
)

// Core provides all functionality in a single interface that can be
// implemented through any means necessary
type Core interface {
	TranslateCore
	RouterCore
	RenderCore
	BindCore
	ErrorCore
	SessionCore
}

// NewCore composes a full core from its partial cores
func NewCore(
	roc RouterCore,
	rec RenderCore,
	bc BindCore,
	sc SessionCore,
	tc TranslateCore,
	ec ErrorCore,
) Core {
	return struct {
		RouterCore
		RenderCore
		BindCore
		SessionCore
		TranslateCore
		ErrorCore
	}{
		roc,
		rec,
		bc,
		sc,
		tc,
		ec}
}

// RouterCore provides part of the core that depends on a router.
type RouterCore interface {
	URL(w http.ResponseWriter, r *http.Request, s string, opts ...URLOption) string
	Params(w http.ResponseWriter, r *http.Request) map[string]string
	Route(w http.ResponseWriter, r *http.Request) string
}

// RenderCore provides part of the core that is responsible for rendering responses. It
// Must include a shared directory of templates on which to perform lookups.
type RenderCore interface {
	Render(w http.ResponseWriter, r *http.Request, rr Render, opts ...RenderOption) error
	Templates
}

// BindCore implements part of the core that does request binding. It must include a shared
// decoder for handling url encoded values.
type BindCore interface {
	Bind(w http.ResponseWriter, r *http.Request, b Bind, opts ...BindOption) (bool, error)
	ValuesDecoder
	StructValidator
}

// SessionCore provides part of the core responsible for sessions
type SessionCore interface {
	Session(w http.ResponseWriter, r *http.Request, opts ...SessionOption) Session
	SessionStore
}

// ErrorCore provides part of the core that is responsible for error handling
type ErrorCore interface {
	HandleError(w http.ResponseWriter, r *http.Request, err error)
}

// TranslateCore provides part of the core that is responsible for translation
type TranslateCore interface {
	Translate(w http.ResponseWriter, r *http.Request, m string, opts ...TranslateOption) string
}

// StructValidator implements struct validation
type StructValidator interface {
	ValidateStruct(context.Context, interface{}) error
}

// ValuesDecoder provides url decoding functionality
type ValuesDecoder interface {
	DecodeValues(v interface{}, values url.Values) (err error)
}

// Templates should provide executable templates by name
type Templates interface {
	Lookup(n string) (TemplateExecuter, error)
}

// TemplateExecuter is the interface that must be implemented by template executers
type TemplateExecuter interface {
	Execute(w io.Writer, vars map[string]reflect.Value, data interface{}) (err error)
}

// SessionStore allows for saving and retrieving sessions
type SessionStore interface {
	LoadSession(w http.ResponseWriter, r *http.Request, name string) (Session, error)
}
