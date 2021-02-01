package rb

import (
	"net/http"
)

type redirectRender struct{ loc string }

func Redirect(loc string) Render { return redirectRender{loc} }

func (r redirectRender) RenderHeader(a *App, w http.ResponseWriter, req *http.Request, status int) (int, error) {
	if status < 1 { // no explicit status code set
		status = http.StatusFound
	}

	w.Header().Set("Location", r.loc)
	return status, nil
}

func (r redirectRender) Value() interface{} { return nil }
func (r redirectRender) Render(a *App, wr http.ResponseWriter, req *http.Request) error {
	return nil
}
