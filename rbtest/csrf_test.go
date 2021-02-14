package rbtest

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/csrf"
)

func TestGenerateCSRF(t *testing.T) {
	var k [32]byte
	c, tok := GenerateCSRF(t, k[:])
	if len(tok) != 88 {
		t.Fatalf("got: %v", tok)
	}

	if len(c) != 144 {
		t.Fatalf("got: %v", c)
	}

	w, r := httptest.NewRecorder(), httptest.NewRequest("POST", "/", nil)
	r.Header.Set("X-CSRF-Token", tok)
	r.Header.Set("Cookie", "_gorilla_csrf="+c)

	csrf.Protect(k[:])(http.HandlerFunc(func(http.ResponseWriter, *http.Request) {})).ServeHTTP(w, r)
	if w.Code != 200 {
		t.Fatalf("got: %v", w.Code)
	}
}
