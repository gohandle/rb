package rb_test

import (
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/go-playground/form/v4"
	"github.com/go-playground/validator/v10"
	rb "github.com/gohandle/rb/rbv2"
	"github.com/gohandle/rb/rbv2/rbplayground"
)

func TestDefaultBindCore(t *testing.T) {
	bc := rb.NewBindCore(
		rbplayground.AdaptDecoder(form.NewDecoder()),
		rbplayground.AdaptValidator(validator.New()))

	b := strings.NewReader(url.Values{"foo": {"bar"}}.Encode())
	w, r := httptest.NewRecorder(), httptest.NewRequest("POST", "/", b)
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	var v struct {
		Foo string `form:"foo" validate:"required"`
	}

	if submit, err := bc.Bind(w, r, rb.Form(&v), rb.AndValidate()); !submit || err != nil {
		t.Fatalf("got: %v %v", submit, err)
	}

	if v.Foo != "bar" {
		t.Fatalf("got: %v", v.Foo)
	}

	t.Run("validation fail", func(t *testing.T) {
		var x struct {
			Bar string `validate:"required"`
		}

		w, r := httptest.NewRecorder(), httptest.NewRequest("POST", "/", nil)
		_, err := bc.Bind(w, r, rb.Form(&x), rb.AndValidate())

		if !strings.Contains(err.Error(), "validation") {
			t.Fatalf("got: %v", err)
		}
	})
}
