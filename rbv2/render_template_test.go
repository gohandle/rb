package rb_test

import (
	"net/http/httptest"
	"testing"

	"github.com/CloudyKit/jet/v6"
	rb "github.com/gohandle/rb/rbv2"
	"github.com/gohandle/rb/rbv2/rbjet"
)

func TestRenderTemplate(t *testing.T) {
	l := jet.NewInMemLoader()
	l.Set("foo.html", `foo, {{foo}}`)

	rc := rb.NewRenderCore(rbjet.Adapt(jet.NewSet(l), nil, nil, nil, nil, nil))

	w, r := httptest.NewRecorder(), httptest.NewRequest("GET", "/", nil)
	rc.Render(w, r, rb.Template("foo.html", nil, rb.TemplateVar("foo", "bar")), rb.Status(201))

	if w.Code != 201 {
		t.Fatalf("got: %v", w.Code)
	}

	if w.Body.String() != "foo, bar" {
		t.Fatalf("got: %v", w.Body.String())
	}
}
