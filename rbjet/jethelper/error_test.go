package jethelper_test

import (
	"errors"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/CloudyKit/jet/v6"
	"github.com/go-playground/validator/v10"
	"github.com/gohandle/rb/rbjet"
	"github.com/gohandle/rb/rbjet/jethelper"
)

func TestFieldErrHelper(t *testing.T) {
	l := jet.NewInMemLoader()
	l.Set("foo.html", `{{ rb_field_error(err1, "Foo") }}{{ rb_non_field_error(err1) }}{{ rb_field_error(err2, "Foo") }}{{ rb_non_field_error(err2) }}`)
	err1 := errors.New("non-field")
	err2 := validator.New().Struct(struct {
		Foo string `validate:"required"`
	}{})

	tmpl, _ := rbjet.Adapt(jet.NewSet(l), nil, nil, nil, nil, nil,
		jethelper.NewFieldError(), jethelper.NewNonFieldError(), nil, nil).Lookup("foo.html")

	vars := jet.VarMap{}
	vars.Set("err1", err1)
	vars.Set("err2", err2)

	w, r := httptest.NewRecorder(), httptest.NewRequest("GET", "/", nil)
	if err := tmpl.Execute(w, r, vars, nil); err != nil {
		t.Fatalf("got: %v", err)
	}

	if act := w.Body.String(); !strings.Contains(act, "non-field") {
		t.Fatalf("got: %v", act)
	}

	if act := w.Body.String(); !strings.Contains(act, "required") {
		t.Fatalf("got: %v", act)
	}
}
