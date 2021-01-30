package rb

import (
	"github.com/CloudyKit/jet/v6"
	"github.com/go-playground/form/v4"
	"github.com/go-playground/validator/v10"
)

type App struct {
	fdec *form.Decoder
	view *jet.Set
	val  *validator.Validate
}

func New(fdec *form.Decoder, view *jet.Set, val *validator.Validate) *App {
	return &App{fdec, view, val}
}

type RenderBind interface {
	Render
	Bind
}
