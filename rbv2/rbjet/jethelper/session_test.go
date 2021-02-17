package jethelper_test

import (
	"net/http/httptest"
	"testing"

	"github.com/CloudyKit/jet/v6"
	rb "github.com/gohandle/rb/rbv2"
	"github.com/gohandle/rb/rbv2/rbgorilla"
	"github.com/gohandle/rb/rbv2/rbjet"
	"github.com/gohandle/rb/rbv2/rbjet/jethelper"
	"github.com/gorilla/sessions"
)

func TestSessionHelper(t *testing.T) {
	l := jet.NewInMemLoader()
	l.Set("foo.html", `{{ rb_session("foo") }}`)

	sc := rb.NewSessionCore(rbgorilla.AdaptSessionStore(sessions.NewCookieStore(make([]byte, 32))))
	tmpl, _ := rbjet.Adapt(jet.NewSet(l), nil, nil, nil, nil, jethelper.NewSession(sc)).Lookup("foo.html")

	w, r := httptest.NewRecorder(), httptest.NewRequest("GET", "/", nil)
	sc.Session(w, r).Set("foo", "bar")

	if err := tmpl.Execute(w, r, nil, nil); err != nil {
		t.Fatalf("got: %v", err)
	}

	if act := w.Body.String(); act != "bar" {
		t.Fatalf("got: %v", act)
	}
}
