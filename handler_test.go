package rb_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gohandle/rb"
	"go.uber.org/zap"
)

func TestHandler(t *testing.T) {
	app := rb.New(zap.NewNop(), nil, nil, nil, nil, nil)

	t.Run("ok", func(t *testing.T) {
		h := app.Action(func(w http.ResponseWriter, r *http.Request) error {
			return app.Respond(w, r, rb.Redirect("/foo"))
		})

		w, r := httptest.NewRecorder(), httptest.NewRequest("GET", "/", nil)
		h.ServeHTTP(w, r)

		if w.Code != 302 {
			t.Fatalf("got: %v", w.Code)
		}
	})

	t.Run("error", func(t *testing.T) {
		h := app.Action(func(w http.ResponseWriter, r *http.Request) error {
			return errors.New("foo")
		})

		w, r := httptest.NewRecorder(), httptest.NewRequest("GET", "/", nil)
		h.ServeHTTP(w, r)

		if w.Code != 500 {
			t.Fatalf("got: %v", w.Code)
		}
	})

}
