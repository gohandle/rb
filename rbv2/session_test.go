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
	sc.Session(w, r).Set("foo", "bar")

	s := rbtest.ReadSession(t, sc, rb.DefaultCookieName, w.Header().Get("Set-Cookie"))
	if act := s.Get("foo"); act != "bar" {
		t.Fatalf("got: %v", act)
	}
}
