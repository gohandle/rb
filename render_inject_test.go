package rb

import (
	"net/http/httptest"
	"testing"

	"github.com/CloudyKit/jet/v6"
)

func TestInjector(t *testing.T) {
	a := New(nil, jet.NewSet(jet.NewInMemLoader()), nil, nil, nil)
	t.Run("ok", func(t *testing.T) {
		w, r := httptest.NewRecorder(), httptest.NewRequest("GET", "/", nil)

		var called int
		a.Render(w, r, Inject(Redirect("/"), InjectorFunc(func(a *App, v interface{}) error {
			called++
			return nil
		})), Status(303))

		if w.Code != 303 {
			t.Fatalf("got: %v", w.Code)
		}

		if act := w.Header().Get("Location"); act != "/" {
			t.Fatalf("got: %v", act)
		}

		if called != 1 {
			t.Fatalf("got: %v", called)
		}
	})
}
