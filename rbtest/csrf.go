package rbtest

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/csrf"
)

// GenerateCSRF generates the two parts required to pass CSRF protection. The cookie and the token
// that should be send together
func GenerateCSRF(tb testing.TB, k []byte, s ...string) (cookieValue, token string) {
	cookieName := "_rb_csrf"
	if len(s) > 0 {
		cookieName = s[0]
	}

	w, r := httptest.NewRecorder(), httptest.NewRequest("GET", "/", nil)
	csrf.Protect(k[:], csrf.CookieName(cookieName))(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "%s", csrf.Token(r))
	})).ServeHTTP(w, r)

	c, err := parseCookie(w.Header().Get("Set-Cookie"), cookieName)
	if err != nil {
		fatalf(tb, "failed to parse cookie: %v", err)
	}

	return c.Value, w.Body.String()
}
