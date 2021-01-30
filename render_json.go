package rb

import "net/http"

type jsonRender struct {
	v interface{}
}

func (r jsonRender) Execute(wr http.ResponseWriter, req *http.Request) error { return nil }

func JSON(data interface{}) Render { return jsonRender{data} }
