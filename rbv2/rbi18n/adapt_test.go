package rbi18n_test

import (
	"net/http/httptest"
	"testing"

	rb "github.com/gohandle/rb/rbv2"
	"github.com/gohandle/rb/rbv2/rbi18n"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

func TestTranslation(t *testing.T) {
	b := i18n.NewBundle(language.English)
	tc := rbi18n.Adapt(b)

	t.Run("no such message", func(t *testing.T) {
		w, r := httptest.NewRecorder(), httptest.NewRequest("GET", "/", nil)
		if act := tc.Translate(w, r, "bogus.message"); act != "bogus.message" {
			t.Fatalf("got: %v", act)
		}
	})

	t.Run("valid message", func(t *testing.T) {
		b.AddMessages(language.English, &i18n.Message{
			ID:    "message.foo.one",
			Other: "hello",
		})

		w, r := httptest.NewRecorder(), httptest.NewRequest("GET", "/", nil)
		if act := tc.Translate(w, r, "message.foo.one"); act != "hello" {
			t.Fatalf("got: %v", act)
		}
	})

	t.Run("other language", func(t *testing.T) {
		b.AddMessages(language.Dutch, &i18n.Message{
			ID:    "message.foo.one",
			Other: "hallo",
		})

		w, r := httptest.NewRecorder(), httptest.NewRequest("GET", "/", nil)
		r.Header.Set("Accept-Language", "nl")
		if act := tc.Translate(w, r, "message.foo.one"); act != "hallo" {
			t.Fatalf("got: %v", act)
		}
	})

	t.Run("plural", func(t *testing.T) {
		b.AddMessages(language.English, &i18n.Message{
			ID:    "message.foo.one",
			One:   "1 hello",
			Other: "hello",
		})

		w, r := httptest.NewRecorder(), httptest.NewRequest("GET", "/", nil)
		if act := tc.Translate(w, r, "message.foo.one", rb.PluralCount(1)); act != "1 hello" {
			t.Fatalf("got: %v", act)
		}
	})
}
