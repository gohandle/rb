package rb_test

import (
	"net/http/httptest"
	"testing"

	"github.com/gohandle/rb"
)

func TestRender(t *testing.T) {
	a := new(rb.App)

	t.Run("render template", func(t *testing.T) {
		w, r := httptest.NewRecorder(), httptest.NewRequest("GET", "/", nil)
		a.Render(w, r,
			rb.Template("foo.html", "foo", rb.TemplateVar("bar", "rab")),
			rb.Status(201))
	})

	t.Run("render json", func(t *testing.T) {
		w, r := httptest.NewRecorder(), httptest.NewRequest("GET", "/", nil)
		a.Render(w, r, rb.JSON("foo"))
	})

	t.Run("render redirect", func(t *testing.T) {
		w, r := httptest.NewRecorder(), httptest.NewRequest("GET", "/", nil)
		a.Render(w, r, rb.Redirect("/"), rb.Status(303))
	})
}
