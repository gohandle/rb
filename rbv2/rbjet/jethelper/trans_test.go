package jethelper_test

import (
	"net/http/httptest"
	"testing"

	"github.com/CloudyKit/jet/v6"
	"github.com/gohandle/rb/rbv2/rbi18n"
	"github.com/gohandle/rb/rbv2/rbjet"
	"github.com/gohandle/rb/rbv2/rbjet/jethelper"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

func TestTransHelper(t *testing.T) {
	l := jet.NewInMemLoader()
	l.Set("foo.html", `{{ t("my.msg.id") }}`)

	b := i18n.NewBundle(language.English)
	b.AddMessages(language.English, &i18n.Message{ID: "my.msg.id", Other: "my msg"})

	tc := rbi18n.Adapt(b)
	tmpl, _ := rbjet.Adapt(jet.NewSet(l), nil, nil, nil, jethelper.NewTrans(tc), nil).Lookup("foo.html")

	w, r := httptest.NewRecorder(), httptest.NewRequest("GET", "/", nil)
	if err := tmpl.Execute(w, r, nil, nil); err != nil {
		t.Fatalf("got: %v", err)
	}

	if act := w.Body.String(); act != "my msg" {
		t.Fatalf("got: %v", act)
	}
}
