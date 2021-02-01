package rb_test

import (
	"net/http/httptest"
	"testing"

	"github.com/gohandle/rb"
	"github.com/gorilla/sessions"
	"go.uber.org/zap"
)

func TestSession(t *testing.T) {
	var k [32]byte
	a := rb.New(zap.NewNop(), nil, nil, nil, sessions.NewCookieStore(k[:]), nil)

	t.Run("session setting", func(t *testing.T) {
		w, r := httptest.NewRecorder(), httptest.NewRequest("GET", "/", nil)
		sess := a.Session(w, r, rb.SessionName("foo")).Set("x", "y").Set("z", "b")

		if act := sess.Pop("z").(string); act != "b" {
			t.Fatalf("got: %v", act)
		}

		if act := sess.Get("x").(string); act != "y" {
			t.Fatalf("got: %v", act)
		}

		if act := sess.Pop("z"); act != nil {
			t.Fatalf("got: %v", act)
		}

		if act := w.Header().Get("Set-Cookie"); len(act) < 50 || act[:3] != "foo" {
			t.Fatalf("got: %v", act)
		}
	})
}
