package rb_test

import (
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/go-playground/form/v4"
	"github.com/gohandle/rb"
)

func TestBind(t *testing.T) {
	a := rb.New(form.NewDecoder(), nil)

	t.Run("bind form", func(t *testing.T) {
		b := strings.NewReader((url.Values{"Foo": {"bar"}}).Encode())
		r := httptest.NewRequest("POST", "/", b)
		r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		var v struct{ Foo string }
		a.Bind(r, rb.FormBind(&v))

		if v.Foo != "bar" {
			t.Fatalf("got: %v", v)
		}
	})

	t.Run("bind json", func(t *testing.T) {
		r := httptest.NewRequest("POST", "/", strings.NewReader(`"foo"`))
		var v string
		a.Bind(r, rb.JSON(&v))

		if v != "foo" {
			t.Fatalf("got: %v", v)
		}
	})
}
