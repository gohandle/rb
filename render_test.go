package rb_test

import (
	"net/http/httptest"
	"testing"

	"github.com/CloudyKit/jet/v6"
	"github.com/gohandle/rb"
)

func TestRender(t *testing.T) {
	l := jet.NewInMemLoader()
	l.Set("foo.html", `{{.}}{{ bar }}`)
	a := rb.New(nil, jet.NewSet(l), nil)

	t.Run("render template", func(t *testing.T) {
		w, r := httptest.NewRecorder(), httptest.NewRequest("GET", "/", nil)
		a.Render(w, r,
			rb.Template("foo.html", "foo", rb.TemplateVar("bar", "rab")),
			rb.Status(201))

		if w.Code != 201 {
			t.Fatalf("got: %v", w.Code)
		}

		if w.Body.String() != "foorab" {
			t.Fatalf("got: %v", w.Body.String())
		}
	})

	t.Run("render json", func(t *testing.T) {
		w, r := httptest.NewRecorder(), httptest.NewRequest("GET", "/", nil)
		a.Render(w, r, rb.JSON("foo"))

		if w.Body.String() != `"foo"`+"\n" {
			t.Fatalf("got: %v", w.Body.String())
		}
	})

	t.Run("render redirect", func(t *testing.T) {
		w, r := httptest.NewRecorder(), httptest.NewRequest("GET", "/", nil)
		a.Render(w, r, rb.Redirect("/"), rb.Status(303))

		if w.Code != 303 {
			t.Fatalf("got: %v", w.Code)
		}

		if act := w.Header().Get("Location"); act != "/" {
			t.Fatalf("got: %v", act)
		}
	})
}
