package jethelper_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/CloudyKit/jet/v6"
	"github.com/gohandle/rb/rbgorilla"
	"github.com/gohandle/rb/rbjet"
	"github.com/gohandle/rb/rbjet/jethelper"
	"github.com/gorilla/mux"
)

func TestRouteHelper(t *testing.T) {
	l := jet.NewInMemLoader()
	l.Set("foo.html", `{{ rb_route() }}`)

	m := mux.NewRouter()
	rc := rbgorilla.AdaptRouter(m)
	tmpl, _ := rbjet.Adapt(jet.NewSet(l)).AddHelper(jethelper.NewRoute(rc)).Lookup("foo.html")

	m.Name("foo").Path("/x/{id}/y").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := tmpl.Execute(w, r, nil, nil); err != nil {
			t.Fatalf("got: %v", err)
		}
	})

	w, r := httptest.NewRecorder(), httptest.NewRequest("GET", "/x/555/y", nil)
	m.ServeHTTP(w, r)

	if act := w.Body.String(); act != "foo" {
		t.Fatalf("got: %v", act)
	}
}
