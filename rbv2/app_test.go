package rb_test

import (
	"fmt"
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

func HandleFoo(a rb.App) http.Handler {
	type Page struct {
		Foo  string `form:"foo" validate:"required"`
		Msg  string
		Loc  string
		ID   string
		Curr string
	}

	return a.Action(func(c rb.Ctx) error {
		var p Page

		if _, err := c.Bind(rb.Form(&p, rb.FromSource(rb.PostFormSource)), rb.AndValidate()); err != nil {
			return fmt.Errorf("failed to bind: %w", err)
		}

		p.Msg = c.Translate("page.about.title", rb.PluralCount(1))
		p.Loc = c.URL("foo", rb.URLVar("id", "111"))
		p.ID = c.Params()["id"]
		p.Curr = c.Route()

		c.Session(rb.CookieName("my_sess")).Set("foo", "bar")

		return c.Render(rb.Template("foo.html", p), rb.Status(201))
	})
}

func TestFooExample(t *testing.T) {
	m, l, k, b := mux.NewRouter(), jet.NewInMemLoader(), make([]byte, 32), i18n.NewBundle(language.English)
	c := rbcore.NewDefault(m, jet.NewSet(l), form.NewDecoder(), validator.New(), sessions.NewCookieStore(k), b)
	a := rb.New(c)

	l.Set("foo.html", `{{.Msg}}: {{.Foo}}: {{.Loc}}: {{.ID}}: {{.Curr}}`)
	b.AddMessages(language.English, &i18n.Message{
		ID:    "page.about.title",
		Other: "About us",
	})
	b.AddMessages(language.Dutch, &i18n.Message{
		ID:    "page.about.title",
		Other: "Over ons",
		One:   "Over ons 1",
	})

	m.Name("foo").Path("/foo/{id}").Handler(HandleFoo(a))

	d := strings.NewReader(url.Values{"foo": {"rab"}}.Encode())
	w, r := httptest.NewRecorder(), httptest.NewRequest("POST", "/foo/888", d)
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Set("Accept-Language", "nl")

	m.ServeHTTP(w, r)

	if w.Code != 201 {
		t.Fatalf("got: %v", w.Code)
	}

	if act := w.Body.String(); act != "Over ons 1: rab: /foo/111: 888: foo" {
		t.Fatalf("got: %v", act)
	}

	s := rbtest.ReadSession(t, c, "my_sess", w.Header().Get("Set-Cookie"))
	if act := s.Get("foo"); act != "bar" {
		t.Fatalf("got: %v", act)
	}
}
