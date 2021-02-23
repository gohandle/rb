package rbgorilla_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	rb "github.com/gohandle/rb"
	"github.com/gohandle/rb/rbgorilla"
	"github.com/gorilla/mux"
)

func TestGorillaAdapt(t *testing.T) {
	mr := mux.NewRouter()
	rc := rbgorilla.AdaptRouter(mr)

	t.Run("url", func(t *testing.T) {
		mr.Name("foo").Path("/foo/{pid}/bar")
		if act := rc.URL(nil, nil, "foo",
			rb.URLVar("pid", "rab"),
			rb.BasePath("/dar"),
		); act != "/dar/foo/rab/bar" {
			t.Fatalf("got: %v", act)
		}
	})

	t.Run("abs", func(t *testing.T) {
		if act := rc.URL(nil, nil, "/foo"); act != "/foo" {
			t.Fatalf("got: %v", act)
		}
	})

	t.Run("vars, and current url", func(t *testing.T) {
		paramv, route := "", ""
		mr.Name("bar").Path("/f/{x}/r").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			paramv = rc.Param(w, r, "x")
			route = rc.Route(w, r)
		})

		w, r := httptest.NewRecorder(), httptest.NewRequest("GET", "/f/y/r", nil)
		mr.ServeHTTP(w, r)

		if paramv != "y" {
			t.Fatalf("got: %+v", paramv)
		}

		if route != "bar" {
			t.Fatalf("got: %+v", route)
		}
	})
}
