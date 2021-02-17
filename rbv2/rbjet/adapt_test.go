package rbjet_test

import (
	"net/http/httptest"
	"testing"

	"github.com/CloudyKit/jet/v6"
	"github.com/gohandle/rb/rbv2/rbjet"
)

func TestAdaptedExecute(t *testing.T) {
	l := jet.NewInMemLoader()
	l.Set("foo.html", `{{ rb_response.Header()["Foo"][0] }}{{rb_request.URL.Path}}`)

	tmpl, _ := rbjet.Adapt(jet.NewSet(l), nil).Lookup("foo.html")
	w, r := httptest.NewRecorder(), httptest.NewRequest("GET", "/foo", nil)
	w.Header().Set("foo", "bar")

	if err := tmpl.Execute(w, r, nil, nil); err != nil {
		t.Fatalf("got: %v", err)
	}

	if act := w.Body.String(); act != "bar/foo" {
		t.Fatalf("got: %v", act)
	}
}
