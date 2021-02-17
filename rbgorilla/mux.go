package rbgorilla

import (
	"net/http"

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
	route := ar.r.Get(s)
	if route == nil {
		rb.L(r).Error("no route with the given name", zap.String("route_name", s))
		return ""
	}

	var o rb.URLOptions
	for _, opt := range opts {
		opt(&o)
	}

	loc, err := route.URL(o.Pairs...)
	if err != nil {
		rb.L(r).Error("failed to generate url",
			zap.String("route_name", s), zap.Strings("pairs", o.Pairs), zap.Error(err))
		return ""
	}

	return loc.String()
}

func (ar *adaptedRouter) Use(mw func(http.Handler) http.Handler) {
	ar.r.Use(mw)
}

func (ar *adaptedRouter) Params(w http.ResponseWriter, r *http.Request) map[string]string {
	return mux.Vars(r)
}

func (ar *adaptedRouter) Route(w http.ResponseWriter, r *http.Request) string {
	return mux.CurrentRoute(r).GetName()
}
