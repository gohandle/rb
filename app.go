package rb

import (
	"github.com/CloudyKit/jet/v6"
	"github.com/go-playground/form/v4"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/sessions"
)

type App struct {
	fdec *form.Decoder
	view *jet.Set
	val  *validator.Validate
	sess sessions.Store
}

func New(
	fdec *form.Decoder,
	view *jet.Set,
	val *validator.Validate,
	sess sessions.Store,
) *App {
	return &App{
		fdec: fdec,
		view: view,
		val:  val,
		sess: sess,
	}
}

type RenderBind interface {
	Render
	Bind
}
