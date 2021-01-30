package rb

import (
	"github.com/CloudyKit/jet/v6"
	"github.com/go-playground/form/v4"
)

type App struct {
	fdec *form.Decoder
	view *jet.Set
}

func New(fdec *form.Decoder, view *jet.Set) *App { return &App{fdec, view} }

type RenderBind interface {
	Render
	Bind
}
