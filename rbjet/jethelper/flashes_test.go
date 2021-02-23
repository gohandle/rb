package jethelper_test

import (
	"net/http/httptest"
	"testing"

	"github.com/CloudyKit/jet/v6"
	rb "github.com/gohandle/rb"
	"github.com/gohandle/rb/rbjet"
	"github.com/gohandle/rb/rbjet/jethelper"
)

func TestFlashesFHelper(t *testing.T) {
	l := jet.NewInMemLoader()
	l.Set("foo.html", `{{ rb_flashes() }}`)

	tmpl, _ := rbjet.Adapt(jet.NewSet(l)).AddHelper(jethelper.NewFlashes()).Lookup("foo.html")

	w, r := httptest.NewRecorder(), httptest.NewRequest("GET", "/", nil)
	r = r.WithContext(rb.WithFlashMessages(r.Context(), []string{"foo"}))

	if err := tmpl.Execute(w, r, nil, nil); err != nil {
		t.Fatalf("got: %v", err)
	}

	if act := w.Body.String(); act != "[foo]" {
		t.Fatalf("got: %v", act)
	}
}
