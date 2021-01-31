package rb_test

import (
	"net/http/httptest"
	"testing"

	"github.com/CloudyKit/jet/v6"
	"github.com/gohandle/rb"
	"github.com/gorilla/sessions"
)

func TestSession(t *testing.T) {
	var k [32]byte
	a := rb.New(nil, jet.NewSet(jet.NewInMemLoader()), nil, sessions.NewCookieStore(k[:]), nil)

	t.Run("session setting", func(t *testing.T) {
		w, r := httptest.NewRecorder(), httptest.NewRequest("GET", "/", nil)
		sess := a.Session(w, r, rb.SessionName("foo")).Set("x", "y").Set("z", "b")

		if act := sess.Get("x").(string); act != "y" {
			t.Fatalf("got: %v", act)
		}

		if act := w.Header().Get("Set-Cookie"); len(act) < 50 || act[:3] != "foo" {
			t.Fatalf("got: %v", act)
		}
	})
}
