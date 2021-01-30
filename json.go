package rb

import (
	"encoding/json"
	"net/http"
)

type jsonRenderBind struct {
	v interface{}
}

func (rb jsonRenderBind) Render(a *App, wr http.ResponseWriter, req *http.Request) error {
	return json.NewEncoder(wr).Encode(rb.v)
}

func (rb jsonRenderBind) Value() interface{} { return rb.v }

func (rb jsonRenderBind) Bind(a *App, req *http.Request) error {
	return json.NewDecoder(req.Body).Decode(rb.v)
}

func JSON(data interface{}) RenderBind {
	return jsonRenderBind{data}
}
