package rbtest

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	rb "github.com/gohandle/rb"
	"go.uber.org/zap"
)

// GenerateSession generates a valid cookie value and csrf token for testing requests and optionally
// allows setting session values in pairs
func GenerateSession(tb testing.TB, k []byte, sc rb.SessionCore, pairs ...interface{}) (c *http.Cookie, tok string) {
	w, r := httptest.NewRecorder(), httptest.NewRequest("GET", "/", nil)
	r = r.WithContext(rb.WithRequestLogger(r.Context(), zap.NewNop()))

	rb.NewCSRFMiddlware(sc, rb.BasicErrorHandler)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "%s", rb.CSRFToken(r.Context()))
		for i := 0; i < len(pairs); i += 2 {
			sc.Session(w, r).Set(pairs[i], pairs[i+1])
		}

		if err := sc.SaveSession(w, r, sc.Session(w, r)); err != nil {
			tb.Fatalf("failed to save cookie during generation: %v", err)
		}
	})).ServeHTTP(w, r)

	c, err := parseCookie(w.Header().Get("Set-Cookie"), rb.DefaultSessionName)
	if err != nil {
		tb.Fatalf("failed to parse cookie: %v", err)
	}

	return c, w.Body.String()
}

// ReadSession allows for asserting sessions in tests
func ReadSession(tb testing.TB, sc rb.SessionCore, name, rawCookie string) rb.SessionReader {
	c, err := parseCookie(rawCookie, name)
	if err != nil {
		tb.Fatalf("failed to parse raw cookie: %v", err)
	}

	r := httptest.NewRequest("GET", "/", nil)
	r.AddCookie(c)

	s, err := sc.LoadSession(nil, r, name)
	if err != nil {
		tb.Fatalf("failed to load session: %v", err)
	}

	return s
}

func parseCookie(rawCookies, name string) (*http.Cookie, error) {
	header := http.Header{}
	header.Add("Cookie", rawCookies)
	r := http.Request{Header: header}
	return r.Cookie(name)
}
