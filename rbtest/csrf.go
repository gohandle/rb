package rbtest

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	rb "github.com/gohandle/rb"
	"github.com/gohandle/rb/rbgorilla"
	"github.com/gorilla/sessions"
)

// GenerateCSRF generates a valid cookie value and csrf token for testing requests
func GenerateCSRF(tb testing.TB, k []byte) (c *http.Cookie, tok string) {
	sc := rb.NewSessionCore(rbgorilla.AdaptSessionStore(sessions.NewCookieStore(k)))

	w, r := httptest.NewRecorder(), httptest.NewRequest("GET", "/", nil)
	rb.NewCSRFMiddlware(sc, rb.BasicErrorHandler)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "%s", rb.CSRFToken(r.Context()))
		sc.SaveSession(w, r, sc.Session(w, r))
	})).ServeHTTP(w, r)

	c, err := parseCookie(w.Header().Get("Set-Cookie"), rb.DefaultSessionName)
	if err != nil {
		tb.Fatalf("failed to parse cookie: %v", err)
	}

	return c, w.Body.String()

}
