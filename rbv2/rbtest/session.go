package rbtest

import (
	"net/http"
	"net/http/httptest"
	"testing"

	rb "github.com/gohandle/rb/rbv2"
)

// sessionReader implements rb.Session but only allows reading
// type sessionReader struct{ vals map[interface{}]interface{} }

// func (s *sessionReader) Get(k interface{}) (v interface{}) {
// 	if s.vals == nil {
// 		s.vals = make(map[interface{}]interface{})
// 	}

// 	v, _ = s.vals[k]
// 	return
// }

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
