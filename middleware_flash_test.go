package rb_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	rb "github.com/gohandle/rb"
	"github.com/gohandle/rb/rbgorilla"
	"github.com/gohandle/rb/rbtest"
	"github.com/gorilla/sessions"
)

type testSessCtx struct {
	w  http.ResponseWriter
	r  *http.Request
	sc rb.SessionCore
}

func (c testSessCtx) Request() *http.Request { return c.r }

func (c testSessCtx) Session(o ...rb.SessionOption) rb.Session {
	return c.sc.Session(c.w, c.r)
}

func TestFlashMiddleware(t *testing.T) {
	sc := rb.NewSessionCore(rbgorilla.AdaptSessionStore(
		sessions.NewCookieStore(make([]byte, 32)),
	))

	w1, r1 := httptest.NewRecorder(), httptest.NewRequest("GET", "/", nil)
	t.Run("set flashes", func(t *testing.T) {
		rb.NewFlashMiddleware(sc)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			c := testSessCtx{w, r, sc}
			rb.Flash(c, "foo")
			rb.Flash(c, "bar", "rab")

			if err := sc.SaveSession(w, r, sc.Session(w, r)); err != nil {
				t.Fatalf("got: %v", err)
			}
		})).ServeHTTP(w1, r1)

		s := rbtest.ReadSession(t, sc, rb.DefaultSessionName, w1.Header().Get("Set-Cookie"))

		if act := s.Get(rb.FlashSessionField); strings.Join(act.([]string), ",") != "foo,bar,rab" {
			t.Fatalf("got: %v", act)
		}
	})

	w2, r2 := httptest.NewRecorder(), httptest.NewRequest("GET", "/", nil)
	r2.Header.Set("Cookie", w1.Header().Get("Set-Cookie"))
	t.Run("read flashes", func(t *testing.T) {

		rb.NewFlashMiddleware(sc)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			c := testSessCtx{w, r, sc}
			fmt.Fprintf(w, "%s", rb.Flashes(c))
		})).ServeHTTP(w2, r2)

		if act := w2.Body.String(); act != `[foo bar rab]` {
			t.Fatalf("got: %v", act)
		}
	})

	t.Run("flashes should now be empty", func(t *testing.T) {
		w3, r3 := httptest.NewRecorder(), httptest.NewRequest("GET", "/", nil)
		r3.Header.Set("Cookie", w2.Header().Get("Set-Cookie"))

		rb.NewFlashMiddleware(sc)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, "%s", rb.FlashMessages(r.Context()))
		})).ServeHTTP(w3, r3)

		if act := w3.Body.String(); act != `[]` {
			t.Fatalf("got: %v", act)
		}
	})

}
