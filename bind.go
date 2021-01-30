package rb

import (
	"net/http"
)

type Bind interface {
	Bind(a *App, r *http.Request) error
}

type formBind struct{ v interface{} }

func FormBind(v interface{}) Bind { return formBind{v} }

func (b formBind) Bind(a *App, r *http.Request) error {
	if err := r.ParseForm(); err != nil {
		return err
	}

	return a.fdec.Decode(b.v, r.PostForm)
}

func (a *App) Bind(r *http.Request, b Bind) error {
	return b.Bind(a, r)
}
