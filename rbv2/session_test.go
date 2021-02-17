package rb_test

import (
	"net/http/httptest"
	"testing"

	rb "github.com/gohandle/rb/rbv2"
	"github.com/gohandle/rb/rbv2/rbgorilla"
	"github.com/gohandle/rb/rbv2/rbtest"
	"github.com/gorilla/sessions"
)

func TestSessions(t *testing.T) {
	var k [32]byte
	sc := rb.NewSessionCore(rbgorilla.AdaptSessionStore(sessions.NewCookieStore(k[:])))

	w, r := httptest.NewRecorder(), httptest.NewRequest("GET", "/", nil)
	s := sc.Session(w, r).Set("foo", "bar")

	if err := sc.SaveSession(w, r, s); err != nil {
		t.Fatalf("got: %v", err)
	}

	sr := rbtest.ReadSession(t, sc, rb.DefaultSessionName, w.Header().Get("Set-Cookie"))
	if act := sr.Get("foo"); act != "bar" {
		t.Fatalf("got: %v", act)
	}
}
