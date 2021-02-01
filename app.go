package rb

import (
	"net/http"

	"github.com/CloudyKit/jet/v6"
	"github.com/go-playground/form/v4"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"go.uber.org/zap"
)

type App struct {
	logs *zap.Logger
	fdec *form.Decoder
	view *jet.Set
	val  *validator.Validate
	sess sessions.Store
	mux  *mux.Router

	// ErrHandler can be configured to get called when an error occured during rendering
	ErrHandler func(a *App, w http.ResponseWriter, r *http.Request, err error) error
}

func New(
	logs *zap.Logger,
	fdec *form.Decoder,
	view *jet.Set,
	val *validator.Validate,
	sess sessions.Store,
	mux *mux.Router,
) *App {
	a := &App{
		logs: logs,
		fdec: fdec,
		view: view,
		val:  val,
		sess: sess,
		mux:  mux,
	}

	if view != nil {
		view.AddGlobalFunc("url", a.urlHelper)
		view.AddGlobalFunc("field_error", a.fieldErrorHelper)
		view.AddGlobalFunc("non_field_error", a.nonFieldErrorHelper)
	}

	return a
}

type RenderBind interface {
	Render
	Bind
}
