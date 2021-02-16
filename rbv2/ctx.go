package rb

import (
	"context"
	"net/http"
)

// Ctx provides request-scoped functionality
type Ctx interface {
	context.Context
	Render(rr Render, o ...RenderOption) error
	Bind(b Bind, o ...BindOption) (bool, error)
	Session(o ...SessionOption) Session
	URL(s string, o ...URLOption) string
	Translate(m string, o ...TranslateOption) string
	Request() *http.Request
	Params() map[string]string
	Route() string
}

// NewCtx creates a context that will implement the request-scoped
// functionality by calling the core
func NewCtx(w http.ResponseWriter, r *http.Request, core Core) Ctx {
	return &ctx{r.Context(), w, r, core}
}

// ctx is a simple wrapper around a core that provides the request-scoped
// functionality
type ctx struct {
	context.Context
	wr   http.ResponseWriter
	req  *http.Request
	core Core
}

func (c *ctx) Request() *http.Request { return c.req }

func (c *ctx) Render(rr Render, o ...RenderOption) error {
	return c.core.Render(c.wr, c.Request(), rr, o...)
}

func (c *ctx) Bind(b Bind, o ...BindOption) (bool, error) {
	return c.core.Bind(c.wr, c.Request(), b, o...)
}

func (c *ctx) Session(o ...SessionOption) Session {
	return c.core.Session(c.wr, c.Request(), o...)
}

func (c *ctx) URL(s string, o ...URLOption) string {
	return c.core.URL(c.wr, c.Request(), s, o...)
}

func (c *ctx) Translate(m string, o ...TranslateOption) string {
	return c.core.Translate(c.wr, c.Request(), m, o...)
}

func (c *ctx) Params() map[string]string {
	return c.core.Params(c.wr, c.Request())
}

func (c *ctx) Route() string {
	return c.core.Route(c.wr, c.Request())
}
