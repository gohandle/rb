package rb

import "net/http"

type redirectRender struct{ loc string }

func Redirect(loc string) Render { return redirectRender{} }

func (r redirectRender) Execute(wr http.ResponseWriter, req *http.Request) error {
	return nil
}
