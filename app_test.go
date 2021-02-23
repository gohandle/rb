package rb_test

import (
	"fmt"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/CloudyKit/jet/v6"
	"github.com/go-playground/form/v4"
	"github.com/go-playground/validator/v10"
	rb "github.com/gohandle/rb"
	"github.com/gohandle/rb/rbapp"
	"github.com/gohandle/rb/rbcore"
	"github.com/gohandle/rb/rbtest"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest/observer"
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
		p.ID = c.Param("id")
		p.Curr = c.Route()

		c.Session().Set("foo", "bar")
		rb.Flash(c, "flash!")

		return c.Render(rb.Template("foo.html", p), rb.Status(201))
	})
}

func TestFooExample(t *testing.T) {
	rnd := rand.New(rand.NewSource(1))
	rb.RandRead = rnd.Read

	m, l, k, b := mux.NewRouter(), jet.NewInMemLoader(), make([]byte, 32), i18n.NewBundle(language.English)
	c := rbcore.NewDefault(m, jet.NewSet(l), form.NewDecoder(), validator.New(), sessions.NewCookieStore(k), b)
	zc, obs := observer.New(zap.DebugLevel)
	a := rbapp.New(c, zap.New(zc))

	l.Set("foo.html", `{{.Msg}}: {{.Foo}}: {{.Loc}}: {{.ID}}: {{.Curr}}{{ rb_csrf() }}{{rb_flashes()}}`)
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

	// setup the request, valid csrf and form submission
	cookie, tok := rbtest.GenerateSession(t, k, c, rb.FlashSessionField, []string{"flash!"})
	d := strings.NewReader(url.Values{"foo": {"rab"}, rb.CSRFFormFieldName: {tok}}.Encode())
	w, r := httptest.NewRecorder(), httptest.NewRequest("POST", "/foo/888", d)
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Set("Accept-Language", "nl")
	r.Header.Set("X-Request-ID", "my-req-id")
	r.AddCookie(cookie)

	// serve and assert the response
	m.ServeHTTP(w, r)
	if w.Code != 201 {
		t.Fatalf("got: %v %v", w.Code, w.Body.String())
	}

	if act := w.Body.String(); act != "Over ons 1: rab: /foo/111: 888: foo650YpEeEBF2H88Z88idG6ZWvWiU2eVG6ov9s1HHEg/G5YOSjZgZhEpHMmXNoRVubAMmdaCZ6LffZRGjToCZFuA==[flash!]" {
		t.Fatalf("got: %v", act)
	}

	// assert the session
	s := rbtest.ReadSession(t, c, rb.DefaultSessionName, w.Header().Get("Set-Cookie"))
	if act := s.Get("foo"); act != "bar" {
		t.Fatalf("got: %v", act)
	}

	if act := s.Get(rb.FlashSessionField); act.([]string)[0] != "flash!" {
		t.Fatalf("got: %v", act)
	}

	// assert that logging includes a request id
	var reqid string
	for _, f := range obs.FilterMessage("render complete").All()[0].Context {
		if f.Key == "request_id" {
			reqid = f.String
		}
	}

	if reqid != "my-req-id" {
		t.Fatalf("got: %v", obs)
	}
}
