package rbgorilla

import (
	"net/http"
	"path"

	rb "github.com/gohandle/rb"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

// AdaptRouter will turn a GorillaMux router into a router core that
// can be used for the router part of the rb core
func AdaptRouter(m *mux.Router) rb.RouterCore {
	return &adaptedRouter{m}
}

type adaptedRouter struct {
	r *mux.Router
}

func (ar *adaptedRouter) URL(w http.ResponseWriter, r *http.Request, s string, opts ...rb.URLOption) string {
	var o rb.URLOptions
	for _, opt := range opts {
		opt(&o)
	}

	if len(s) > 0 && s[0] == '/' {
		return path.Join(o.BasePath, s)
	}

	route := ar.r.Get(s)
	if route == nil {
		rb.L(r).Error("no route with the given name", zap.String("route_name", s))
		return ""
	}

	loc, err := route.URL(o.Pairs...)
	if err != nil {
		rb.L(r).Error("failed to generate url",
			zap.String("route_name", s), zap.Strings("pairs", o.Pairs), zap.Error(err))
		return ""
	}

	s = loc.String()
	if o.BasePath != "" {
		s = path.Join(o.BasePath, s)
	}

	return s
}

func (ar *adaptedRouter) Use(mw func(http.Handler) http.Handler) {
	ar.r.Use(mw)
}

func (ar *adaptedRouter) Param(w http.ResponseWriter, r *http.Request, k string) (v string) {
	v, _ = mux.Vars(r)[k]
	return
}

func (ar *adaptedRouter) Route(w http.ResponseWriter, r *http.Request) string {
	return mux.CurrentRoute(r).GetName()
}
