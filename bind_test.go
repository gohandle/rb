package rb_test

import (
	"errors"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/CloudyKit/jet/v6"
	"github.com/go-playground/form/v4"
	"github.com/go-playground/validator/v10"
	"github.com/gohandle/rb"
)

func TestBind(t *testing.T) {
	a := rb.New(form.NewDecoder(), jet.NewSet(jet.NewInMemLoader()), validator.New(), nil, nil)

	t.Run("bind form", func(t *testing.T) {
		b := strings.NewReader((url.Values{"Foo": {"bar"}}).Encode())
		r := httptest.NewRequest("POST", "/", b)
		r.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		var v struct{ Foo string }
		if err := a.Bind(r, rb.FormBind(&v)); err != nil {
			t.Fatalf("got: %v", err)
		}

		if v.Foo != "bar" {
			t.Fatalf("got: %v", v)
		}
	})

	t.Run("bind json", func(t *testing.T) {
		r := httptest.NewRequest("POST", "/", strings.NewReader(`"foo"`))
		var v string
		if err := a.Bind(r, rb.JSON(&v)); err != nil {
			t.Fatalf("got: %v", err)
		}

		if v != "foo" {
			t.Fatalf("got: %v", v)
		}
	})

	t.Run("bind json and validate fail", func(t *testing.T) {
		r := httptest.NewRequest("POST", "/", strings.NewReader(`{"foo": ""}`))
		var v struct {
			Foo string `validate:"required"`
		}

		err := a.Bind(r, rb.JSON(&v), rb.AndValidate())

		var verr validator.ValidationErrors
		if !errors.As(err, &verr) {
			t.Fatalf("got: %v", err)
		}

		if len(verr) != 1 || verr[0].Field() != "Foo" {
			t.Fatalf("got: %v", verr)
		}
	})
}
