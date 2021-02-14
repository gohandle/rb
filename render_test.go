package rb_test

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/CloudyKit/jet/v6"
	"github.com/go-playground/form/v4"
	"github.com/gohandle/rb"
	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
)

func TestRender(t *testing.T) {
	l := jet.NewInMemLoader()
	l.Set("foo.html", `{{.}}{{ bar }}`)
	a := rb.New(zap.NewNop(), nil, jet.NewSet(l), nil, nil, nil)

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

	t.Run("render redirect status without explit status", func(t *testing.T) {
		w, r := httptest.NewRecorder(), httptest.NewRequest("GET", "/", nil)
		a.Render(w, r, rb.Redirect("/"))

		// NOTE: it is very confusing if the redirec doesn't work. It is usually because not
		// explicit status code is provided. So we'll default
		if w.Code != 302 {
			t.Fatalf("got: %v", w.Code)
		}

		if act := w.Header().Get("Location"); act != "/" {
			t.Fatalf("got: %v", act)
		}
	})
}

type errHeaderRender struct{}

func (r errHeaderRender) RenderHeader(a *rb.App, w http.ResponseWriter, req *http.Request, status int) (int, error) {
	return status, errors.New("expected error")
}

func (r errHeaderRender) String() string     { return "test" }
func (r errHeaderRender) Value() interface{} { return nil }
func (r errHeaderRender) Render(a *rb.App, wr http.ResponseWriter, req *http.Request) error {
	return nil
}

func TestRenderError(t *testing.T) {
	l := jet.NewInMemLoader()
	a := rb.New(zap.NewNop(), nil, jet.NewSet(l), nil, nil, nil)

	t.Run("body render error", func(t *testing.T) {
		w, r := httptest.NewRecorder(), httptest.NewRequest("GET", "/", nil)
		a.Render(w, r, rb.Template("bogus.html", "bogus"), rb.Status(202))

		// NOTE: The error occurs after the header is already written, so we expect the
		// status code to be whateve was configued
		if w.Code != 202 {
			t.Fatalf("got: %v", w.Code)
		}

		if act := w.Body.String(); !strings.Contains(act, "failed to get template") {
			t.Fatalf("got: %v", act)
		}
	})

	t.Run("header render error", func(t *testing.T) {
		w, r := httptest.NewRecorder(), httptest.NewRequest("GET", "/", nil)
		a.Render(w, r, errHeaderRender{}, rb.Status(202))

		// NOTE: here we DO expect the status to be 500, since the header failed writing
		if w.Code != 500 {
			t.Fatalf("got: %v", w.Code)
		}

		if act := w.Body.String(); !strings.Contains(act, "expected error") {
			t.Fatalf("got: %v", act)
		}
	})

	t.Run("template render error halfway", func(t *testing.T) {
		l.Set("fail.html", `foobar {{.bogus}}foo`)
		w, r := httptest.NewRecorder(), httptest.NewRequest("GET", "/", nil)
		a.Render(w, r, rb.Template("fail.html", "fail"), rb.Status(202))

		// NOTE: because the response was already underway when the error occured
		if w.Code != 202 {
			t.Fatalf("got: %v", w.Code)
		}

		// we expect the template error to be rendered as it occurs in the response
		if !strings.HasPrefix(w.Body.String(), `foobar `) {
			t.Fatalf("got: %v", w.Body.String())
		}

		if !strings.HasSuffix(w.Body.String(), `cannot index slice/array/string with type string`+"\n") {
			t.Fatalf("got: %v", w.Body.String())
		}
	})

	t.Run("failing error handling", func(t *testing.T) {
		defer func() {
			var s string
			if r := recover(); r != nil {
				s = fmt.Sprintf("%v", r)
			}

			if s != `rb: failed to handle error, original error: 'expected error', error handler error: expected` {
				t.Fatalf("got: %v", s)
			}
		}()

		a := rb.New(zap.NewNop(), nil, jet.NewSet(l), nil, nil, nil, rb.ErrorHandler(
			func(a *rb.App, w http.ResponseWriter, r *http.Request, err error) error {
				return errors.New("expected")
			}))

		w, r := httptest.NewRecorder(), httptest.NewRequest("GET", "/", nil)
		a.Render(w, r, errHeaderRender{})
	})
}

func TestLogging(t *testing.T) {
	zc, obs := observer.New(zap.DebugLevel)

	templates := jet.NewInMemLoader()
	a := rb.New(zap.New(zc), nil, jet.NewSet(templates), nil, nil, nil)
	inj := func(a *rb.App, w http.ResponseWriter, req *http.Request, v interface{}) error { return nil }

	t.Run("render", func(t *testing.T) {
		w, r := httptest.NewRecorder(), httptest.NewRequest("GET", "/", nil)
		a.Render(w, r, rb.Inject(rb.JSON("foo"), rb.InjectorFunc(inj)))

		if obs.FilterMessage("start render").Len() != 1 {
			t.Fatalf("got: %v", obs.All())
		}

		if obs.FilterMessage("render complete").Len() != 1 {
			t.Fatalf("got: %v", obs.All())
		}
	})
}

func TestRenderLogging(t *testing.T) {
	t.Run("render json", func(t *testing.T) {
		lbuf := bytes.NewBuffer(nil)
		ws := zapcore.AddSync(lbuf)
		zc := zapcore.NewCore(zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig()), ws, zap.DebugLevel)
		a := rb.New(zap.New(zc), nil, nil, nil, nil, nil)

		w, r := httptest.NewRecorder(), httptest.NewRequest("GET", "/", nil)
		a.Render(w, r, rb.JSON("foo"))

		if !strings.Contains(lbuf.String(), "no explicit status provided, set default") {
			t.Fatalf("got: %v", lbuf.String())
		}
	})
}

func TestRenderCSRFInTemplate(t *testing.T) {
	var k [32]byte
	l := jet.NewInMemLoader()
	l.Set("foo.html", `{{ csrf_token }}`)
	m := mux.NewRouter()
	a := rb.New(zap.NewNop(), form.NewDecoder(), jet.NewSet(l), nil, nil, m,
		rb.ProtectFromCSRF(k[:], csrf.CookieName("_my_csrf"), csrf.FieldName("my_csrf_token")))

	type testSubmit struct {
		Foo string `form:"foo"`
	}

	var s testSubmit
	m.Handle("/foo", a.Action(func(w http.ResponseWriter, r *http.Request) error {
		if err := a.Bind(r, rb.Form(&s)); err != nil {
			t.Fatalf("got: %v", err)
		}

		return a.Respond(w, r, rb.Template("foo.html", nil))
	}))

	w, r := httptest.NewRecorder(), httptest.NewRequest("GET", "/foo", nil)
	m.ServeHTTP(w, r)

	cookie, token := w.Header().Get("Set-Cookie"), w.Body.String()
	if !strings.HasPrefix(cookie, "_my_csrf") {
		t.Fatalf("got: %v", cookie)
	}

	if act := len(token); act != 88 {
		t.Fatalf("got: %v", act)
	}

	t.Run("submit with valid token and cookie", func(t *testing.T) {
		b := strings.NewReader(url.Values{
			"my_csrf_token": {strings.TrimSpace(w.Body.String())},
			"foo":           {"bar"},
		}.Encode())

		w, r := httptest.NewRecorder(), httptest.NewRequest("POST", "/foo", b)
		r.Header.Set("Cookie", cookie)
		r.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		m.ServeHTTP(w, r)
		if w.Code != 200 {
			t.Fatalf("got: %v", w.Body.String())
		}

		// form binding should have worked as normal
		if s.Foo != "bar" {
			t.Fatalf("got: %v", s.Foo)
		}
	})

	t.Run("without cookie", func(t *testing.T) {
		b := strings.NewReader(url.Values{
			"my_csrf_token": {strings.TrimSpace(w.Body.String())},
		}.Encode())

		w, r := httptest.NewRecorder(), httptest.NewRequest("POST", "/foo", b)
		r.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		m.ServeHTTP(w, r)
		if w.Code != http.StatusForbidden {
			t.Fatalf("got: %v", w.Body.String())
		}
	})
}
