package rb

import (
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/CloudyKit/jet/v6"
	"github.com/go-playground/form/v4"
	"github.com/go-playground/validator/v10"
)

type testValError struct{}

func (e testValError) Error() string { return "foo_error" }
func (e testValError) Field() string { return "foo" }

func TestValidationRendering(t *testing.T) {
	val, templates := validator.New(), jet.NewInMemLoader()
	a := New(form.NewDecoder(), jet.NewSet(templates), val, nil, nil)

	t.Run("render nil", func(t *testing.T) {
		templates.Set("t1.html", `{{ field_error(., "field_foo") }}`)

		var v error

		w, r := httptest.NewRecorder(), httptest.NewRequest("GET", "/", nil)
		a.Render(w, r, Template("t1.html", v))

		// NOTE: passing in nil directly should also result in nothing being rendered
		if act := w.Body.String(); act != `` {
			t.Fatalf("got: %v", act)
		}
	})

	t.Run("render nil", func(t *testing.T) {
		templates.Set("t2.html", `{{ field_error(.Err, "field_foo") }}`)

		var v struct{ Err error }

		w, r := httptest.NewRecorder(), httptest.NewRequest("GET", "/", nil)
		a.Render(w, r, Template("t2.html", v))

		// NOTE: resolving to a field that is a nil error should also result in nothing
		if act := w.Body.String(); act != `` {
			t.Fatalf("got: %v", act)
		}
	})

	t.Run("ok for custom error", func(t *testing.T) {
		templates.Set("t3.html", `{{ field_error(., "foo") }}`)

		var v error
		v = testValError{}

		w, r := httptest.NewRecorder(), httptest.NewRequest("GET", "/", nil)
		a.Render(w, r, Template("t3.html", v))

		if act := w.Body.String(); act != `foo_error` {
			t.Fatalf("got: %v", act)
		}
	})

	t.Run("ok for bind with validate", func(t *testing.T) {
		templates.Set("t4.html", `{{ field_error(.Err, "Foo") }}`)

		var v struct {
			Err error
			Foo string `form:"foo" validate:"required"`
		}

		b := strings.NewReader((url.Values{"foo": {""}}).Encode())
		w, r := httptest.NewRecorder(), httptest.NewRequest("POST", "/", b)
		r.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		v.Err = a.Bind(r, FormBind(&v), AndValidate())
		if v.Err == nil {
			t.Fatalf("got: %v", v.Err)
		}

		a.Render(w, r, Template("t4.html", v))

		if act := w.Body.String(); act != `failed to validate: Key: &#39;Foo&#39; Error:Field validation for &#39;Foo&#39; failed on the &#39;required&#39; tag` {
			t.Fatalf("got: %v", act)
		}
	})
}
