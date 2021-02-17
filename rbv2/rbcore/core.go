package rbcore

import (
	"github.com/CloudyKit/jet/v6"
	"github.com/go-playground/form/v4"
	"github.com/go-playground/validator/v10"
	rb "github.com/gohandle/rb/rbv2"
	"github.com/gohandle/rb/rbv2/rbgorilla"
	"github.com/gohandle/rb/rbv2/rbi18n"
	"github.com/gohandle/rb/rbv2/rbjet"
	"github.com/gohandle/rb/rbv2/rbjet/jethelper"
	"github.com/gohandle/rb/rbv2/rbplayg"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

// New inits a new core from all its sub cores
func New(
	roc rb.RouterCore,
	rec rb.RenderCore,
	bc rb.BindCore,
	sc rb.SessionCore,
	tc rb.TranslateCore,
	ec rb.ErrorCore,
) rb.Core {
	return struct {
		rb.RouterCore
		rb.RenderCore
		rb.BindCore
		rb.SessionCore
		rb.TranslateCore
		rb.ErrorCore
	}{
		roc,
		rec,
		bc,
		sc,
		tc,
		ec,
	}
}

// NewDefault creates a core with default dependencies
func NewDefault(router *mux.Router, jset *jet.Set, fdec *form.Decoder, val *validator.Validate, ss sessions.Store, bundle *i18n.Bundle) rb.Core {
	rc, tc, sc := rbgorilla.AdaptRouter(router), rbi18n.Adapt(bundle), rb.NewSessionCore(rbgorilla.AdaptSessionStore(ss))
	return New(
		rc,
		rb.NewRenderCore(rbjet.Adapt(jset,
			jethelper.NewURL(rc),
			jethelper.NewParams(rc),
			jethelper.NewRoute(rc),
			jethelper.NewTrans(tc),
			jethelper.NewSession(sc),
		)),
		rb.NewBindCore(rbplayg.AdaptDecoder(fdec), rbplayg.AdaptValidator(val)),
		sc,
		tc,
		rb.BasicErrorHandler,
	)
}
