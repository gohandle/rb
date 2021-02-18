package jethelper_test

import (
	"net/http/httptest"
	"testing"

	"github.com/CloudyKit/jet/v6"
	"github.com/gohandle/rb/rbgorilla"
	"github.com/gohandle/rb/rbjet"
	"github.com/gohandle/rb/rbjet/jethelper"
	"github.com/gorilla/mux"
)

func TestURLHelper(t *testing.T) {
	l := jet.NewInMemLoader()
	l.Set("foo.html", `{{ rb_url("foo", "id", "1234") }}`)

	m := mux.NewRouter()
	m.Name("foo").Path("/x/{id}/y")
	rc := rbgorilla.AdaptRouter(m)

	tmpl, _ := rbjet.Adapt(jet.NewSet(l), jethelper.NewURL(rc), nil, nil, nil, nil, nil, nil, nil, nil).Lookup("foo.html")
	w, r := httptest.NewRecorder(), httptest.NewRequest("GET", "/foo", nil)

	if err := tmpl.Execute(w, r, nil, nil); err != nil {
		t.Fatalf("got: %v", err)
	}

	if act := w.Body.String(); act != "/x/1234/y" {
		t.Fatalf("got: %v", act)
	}
}
