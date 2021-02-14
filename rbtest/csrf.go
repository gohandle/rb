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
func GenerateCSRF(tb testing.TB, k []byte) (cookieValue, token string) {
	w, r := httptest.NewRecorder(), httptest.NewRequest("GET", "/", nil)
	csrf.Protect(k[:])(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "%s", csrf.Token(r))
	})).ServeHTTP(w, r)

	c, err := parseCookie(w.Header().Get("Set-Cookie"), "_gorilla_csrf")
	if err != nil {
		fatalf(tb, "failed to parse cookie: %v", err)
	}

	return c.Value, w.Body.String()
}
