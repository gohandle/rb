package rb_test

import (
	"bytes"
	"errors"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/go-playground/form/v4"
	"github.com/go-playground/validator/v10"
	"github.com/gohandle/rb"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func TestBind(t *testing.T) {
	a := rb.New(zap.NewNop(), form.NewDecoder(), nil, validator.New(), nil, nil)

	t.Run("bind form", func(t *testing.T) {
		b := strings.NewReader((url.Values{"Foo": {"bar"}}).Encode())
		r := httptest.NewRequest("POST", "/?Bar=rab", b)
		r.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		t.Run("default", func(t *testing.T) {
			var v struct {
				Foo string
				Bar string
			}
			if err := a.Bind(r, rb.Form(&v)); err != nil {
				t.Fatalf("got: %v", err)
			}

			if v.Foo != "bar" || v.Bar != "rab" {
				t.Fatalf("got: %v", v)
			}
		})

		t.Run("post form", func(t *testing.T) {
			var v struct {
				Foo string
				Bar string
			}
			if err := a.Bind(r, rb.Form(&v, rb.FromSource(rb.PostFormSource))); err != nil {
				t.Fatalf("got: %v", err)
			}

			if v.Foo != "bar" || v.Bar != "" {
				t.Fatalf("got: %v", v)
			}
		})

		t.Run("query", func(t *testing.T) {
			var v struct {
				Foo string
				Bar string
			}
			if err := a.Bind(r, rb.Form(&v, rb.FromSource(rb.QuerySource))); err != nil {
				t.Fatalf("got: %v", err)
			}

			if v.Foo != "" || v.Bar != "rab" {
				t.Fatalf("got: %v", v)
			}
		})
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

func TestBindLogging(t *testing.T) {
	t.Run("bind form", func(t *testing.T) {
		lbuf := bytes.NewBuffer(nil)
		ws := zapcore.AddSync(lbuf)
		zc := zapcore.NewCore(zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig()), ws, zap.DebugLevel)
		a := rb.New(zap.New(zc), form.NewDecoder(), nil, nil, nil, nil)

		var v struct {
			Foo string
		}

		b := strings.NewReader((url.Values{"Foo": {"bar"}}).Encode())
		r := httptest.NewRequest("POST", "/?Bar=rab", b)
		r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		a.Bind(r, rb.Form(&v))

		if !strings.Contains(lbuf.String(), "bind validate set to false, bind done") {
			t.Fatalf("got: %v", lbuf.String())
		}

		if !strings.Contains(lbuf.String(), "binding from form") {
			t.Fatalf("got: %v", lbuf.String())
		}
	})
}
