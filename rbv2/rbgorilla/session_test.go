package rbgorilla_test

import (
	"net/http/httptest"
	"testing"

	"github.com/gohandle/rb/rbv2/rbgorilla"
	"github.com/gorilla/sessions"
)

func TestSessionAdapt(t *testing.T) {
	var k [32]byte

	ss := rbgorilla.AdaptSessionStore(sessions.NewCookieStore(k[:]))
	w, r := httptest.NewRecorder(), httptest.NewRequest("GET", "/", nil)
	s, err := ss.LoadSession(w, r, "foo")
	if err != nil {
		t.Fatalf("got: %v", err)
	}

	if w.Header().Get("Set-Cookie") != "" {
		t.Fatalf("got: %v", w.Header())
	}

	s.Set("foo", "bar")
	s.Set("rab", "dar")

	if err := ss.SaveSession(w, r, s); err != nil {
		t.Fatalf("got: %v", err)
	}

	if w.Header().Get("Set-Cookie") == "" {
		t.Fatalf("got: %v", w.Header())
	}

	if act := s.Pop("foo"); act != "bar" {
		t.Fatalf("got: %v", act)
	}

	if act := s.Pop("foo"); act != nil {
		t.Fatalf("got: %v", act)
	}

	if act := s.Get("rab"); act != "dar" {
		t.Fatalf("got: %v", act)
	}
}
