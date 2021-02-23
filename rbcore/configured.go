package rbcore

import (
	"net/http"

	"github.com/gohandle/rb"
)

type configured struct {
	rb.Core
	opts Options
}

// Options hold default options that can be configured
type Options struct {
	URLOptions []rb.URLOption
}

// Configured create a core that prepends all variadic configuration slices
// with default options.
func Configured(c rb.Core, opts Options) rb.Core {
	return &configured{Core: c, opts: opts}
}

func (c *configured) URL(w http.ResponseWriter, r *http.Request, s string, opts ...rb.URLOption) string {
	return c.Core.URL(w, r, s, append(c.opts.URLOptions, opts...)...)
}
