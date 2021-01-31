package rbtest

import (
	"encoding/base64"
	"net/http/httptest"
	"testing"

	"github.com/CloudyKit/jet/v6"
	"github.com/gohandle/rb"
	"github.com/gorilla/sessions"
)

func TestSessionReading(t *testing.T) {
	k1, _ := base64.StdEncoding.DecodeString("bDQ1TVQbmuYlaDZp415XGab2Q3xMiLl/wD+Nc+ouy4M=")
	cs := sessions.NewCookieStore(k1)
	a := rb.New(nil, jet.NewSet(jet.NewInMemLoader()), nil, cs, nil)
	w, r := httptest.NewRecorder(), httptest.NewRequest("GET", "/", nil)
	a.Session(w, r).Set("foo", "bar")

	sess := ReadSession(t, cs, "rb", w.Header().Get("Set-Cookie"))
	if act := sess.Get("foo").(string); act != "bar" {
		t.Fatalf("got: %v", act)
	}
}
