package rb_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	rb "github.com/gohandle/rb"
	"github.com/gohandle/rb/rbgorilla"
	"github.com/gohandle/rb/rbtest"
	"github.com/gorilla/sessions"
)

func TestCSRFMiddleware(t *testing.T) {
	ss := sessions.NewCookieStore(make([]byte, 32))
	sc := rb.NewSessionCore(rbgorilla.AdaptSessionStore(ss))

	w1, r1 := httptest.NewRecorder(), httptest.NewRequest("GET", "/", nil)
	rb.NewCSRFMiddlware(sc, rb.BasicErrorHandler)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "%s", rb.CSRFToken(r.Context()))
		if err := sc.SaveSession(w, r, sc.Session(w, r)); err != nil {
			t.Fatalf("failed to save session: %v", err)
		}
	})).ServeHTTP(w1, r1)

	s := rbtest.ReadSession(t, sc, rb.DefaultSessionName, w1.Header().Get("Set-Cookie"))
	if act := w1.Body.String(); len(act) != 88 {
		t.Fatalf("got: %v", act)
	}

	if act, ok := s.Get(rb.CSRFSessionFieldName).([]byte); !ok || len(act) != 32 {
		t.Fatalf("got: %v %v", act, ok)
	}

	t.Run("post", func(t *testing.T) {
		b2 := strings.NewReader(url.Values{
			rb.CSRFFormFieldName: {w1.Body.String()},
		}.Encode())

		w2, r2 := httptest.NewRecorder(), httptest.NewRequest("POST", "/", b2)
		r2.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		r2.Header.Set("Cookie", w1.Header().Get("Set-Cookie"))

		rb.NewCSRFMiddlware(sc, rb.BasicErrorHandler)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, "%s", rb.CSRFToken(r.Context()))
		})).ServeHTTP(w2, r2)

		if w2.Code != 200 {
			t.Fatalf("got: %v", w2.Code)
		}
	})
}
