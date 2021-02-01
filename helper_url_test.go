package rb

import (
	"net/http/httptest"
	"testing"

	"github.com/CloudyKit/jet/v6"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

func TestURLHelper(t *testing.T) {
	router, templates := mux.NewRouter(), jet.NewInMemLoader()
	router.Name("route_b").Path("/b/{id}/b")
	a := New(zap.NewNop(), nil, jet.NewSet(templates), nil, nil, router)

	t.Run("ok", func(t *testing.T) {
		templates.Set("t1.html", `{{ url("route_b", "id", "1234") }}`)

		w, r := httptest.NewRecorder(), httptest.NewRequest("GET", "/a", nil)
		a.Render(w, r, Template("t1.html", "foo"))

		if act := w.Body.String(); act != `/b/1234/b` {
			t.Fatalf("got: %v", act)
		}
	})
}
