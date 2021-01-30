package rb

import "net/http"

type redirectRender struct{ loc string }

func Redirect(loc string) Render { return redirectRender{loc} }

func (r redirectRender) RenderHeader(a *App, w http.ResponseWriter, req *http.Request, status int) error {
	w.Header().Set("Location", r.loc)
	w.WriteHeader(status)
	return nil
}

func (r redirectRender) Render(a *App, wr http.ResponseWriter, req *http.Request) error {
	return nil
}
