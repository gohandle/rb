package rbtest

import (
	"net/http"
	"testing"

	"github.com/gohandle/rb"
	"github.com/gorilla/sessions"
)

// testSession implements rb.Session but only allows reading
type testSession struct{ vals map[interface{}]interface{} }

func (s *testSession) Get(k interface{}) (v interface{}) {
	if s.vals == nil {
		s.vals = make(map[interface{}]interface{})
	}

	v, _ = s.vals[k]
	return
}

func ReadSession(tb testing.TB, s *sessions.CookieStore, name, rawCookie string) rb.SessionReader {
	c, err := parseCookie(rawCookie, name)
	if err != nil {
		fatalf(tb, "failed to parse Set-Cookie header name=%s (%s): %v", name, rawCookie, err)
	}

	sess := sessions.NewSession(s, name)
	if err = s.Codecs[0].Decode(name, c.Value, &sess.Values); err != nil {
		fatalf(tb, "failed to decode cookie: %v", err)
	}

	return &testSession{vals: sess.Values}
}

func parseCookie(rawCookies, name string) (*http.Cookie, error) {
	header := http.Header{}
	header.Add("Cookie", rawCookies)
	r := http.Request{Header: header}
	return r.Cookie(name)
}
