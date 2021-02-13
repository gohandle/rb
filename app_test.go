package rb_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/CloudyKit/jet/v6"
	"github.com/gohandle/rb"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest/observer"
)

type testApp1 struct{ *rb.App }

func (a *testApp1) handleFoo() http.HandlerFunc {
	type page struct {
		Name string
	}

	return func(w http.ResponseWriter, r *http.Request) {
		p, err := page{}, (error)(nil)
		defer a.Render(w, r, rb.Template("foo.html", &p), rb.WithError(&err))

		if fail := r.URL.Query().Get("fail"); fail != "" {
			err = errors.New(fail)
		}

		p.Name = "world"
	}
}

func TestAppWithExplictErrorOptionForRender(t *testing.T) {
	zc, obs := observer.New(zap.DebugLevel)
	v := jet.NewInMemLoader()
	v.Set("foo.html", `hello, {{.Name}}`)

	a := testApp1{rb.New(zap.New(zc), nil, jet.NewSet(v), nil, nil, nil)}
	w, r := httptest.NewRecorder(), httptest.NewRequest("GET", "/", nil)
	a.handleFoo().ServeHTTP(w, r)

	if act := w.Body.String(); act != "hello, world" {
		t.Fatalf("got: %v", act)
	}

	t.Run("or with an error", func(t *testing.T) {
		w, r := httptest.NewRecorder(), httptest.NewRequest("GET", "/?fail=foo", nil)
		a.handleFoo().ServeHTTP(w, r)

		if obs.FilterMessage("explicit render error").Len() != 1 {
			t.Fatalf("got: %v", obs)
		}
	})
}
