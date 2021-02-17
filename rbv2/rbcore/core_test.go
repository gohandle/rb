package rbcore_test

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/CloudyKit/jet/v6"
	"github.com/go-playground/form/v4"
	"github.com/go-playground/validator/v10"
	rb "github.com/gohandle/rb/rbv2"
	"github.com/gohandle/rb/rbv2/rbcore"
	"github.com/gohandle/rb/rbv2/rbtest"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

func TestDefaultCore(t *testing.T) {
	m, l, k, b := mux.NewRouter(), jet.NewInMemLoader(), make([]byte, 32), i18n.NewBundle(language.English)
	c := rbcore.NewDefault(m, jet.NewSet(l), form.NewDecoder(), validator.New(), sessions.NewCookieStore(k), b)

	l.Set("about.html", `{{.Msg}}: {{.Foo}}`)
	b.AddMessages(language.English, &i18n.Message{
		ID:    "page.about.title",
		Other: "About us",
	})
	b.AddMessages(language.Dutch, &i18n.Message{
		ID:    "page.about.title",
		Other: "Over ons",
	})

	m.Name("about").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var v struct {
			Foo string `form:"foo" validate:"required"`
			Msg string
		}

		var err error
		if _, err = c.Bind(w, r, rb.Form(&v), rb.AndValidate()); err != nil {
			t.Fatalf("got: %v", err)
		}

		v.Msg = c.Translate(w, r, "page.about.title")

		s := c.Session(w, r).Set("foo", "bar")

		if err = c.SaveSession(w, r, s); err != nil {
			t.Fatalf("got: %v", err)
		}

		if err = c.Render(w, r, rb.Template("about.html", v)); err != nil {
			t.Fatalf("got: %v", err)
		}
	})

	d := strings.NewReader(url.Values{"foo": {"bar"}}.Encode())
	w, r := httptest.NewRecorder(), httptest.NewRequest("POST", "/", d)
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Set("Accept-Language", "nl")
	m.ServeHTTP(w, r)

	if act := w.Body.String(); act != "Over ons: bar" {
		t.Fatalf("got: %v", act)
	}

	s := rbtest.ReadSession(t, c, rb.DefaultCookieName, w.Header().Get("Set-Cookie"))
	if act := s.Get("foo"); act != "bar" {
		t.Fatalf("got: %v", act)
	}
}
