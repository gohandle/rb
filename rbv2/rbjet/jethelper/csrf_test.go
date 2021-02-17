package jethelper_test

import (
	"net/http/httptest"
	"testing"

	"github.com/CloudyKit/jet/v6"
	"github.com/gohandle/rb/rbv2/rbcsrf"
	"github.com/gohandle/rb/rbv2/rbjet"
	"github.com/gohandle/rb/rbv2/rbjet/jethelper"
)

func TestCSRFHelper(t *testing.T) {
	l := jet.NewInMemLoader()
	l.Set("foo.html", `{{ rb_csrf() }}`)

	tmpl, _ := rbjet.Adapt(jet.NewSet(l),
		nil, nil, nil, nil, nil, nil, nil, jethelper.NewCSRF()).Lookup("foo.html")

	w, r := httptest.NewRecorder(), httptest.NewRequest("GET", "/x/555/y", nil)
	r = r.WithContext(rbcsrf.WithToken(r.Context(), "foo"))

	if err := tmpl.Execute(w, r, nil, nil); err != nil {
		t.Fatalf("got: %v", err)
	}

	if act := w.Body.String(); act != "foo" {
		t.Fatalf("got: %v", act)
	}
}
