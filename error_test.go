package rb_test

import (
	"errors"
	"net/http/httptest"
	"testing"

	rb "github.com/gohandle/rb"
)

func TestCoreError(t *testing.T) {
	ec := rb.BasicErrorHandler
	w, r := httptest.NewRecorder(), httptest.NewRequest("GET", "/", nil)
	ec.HandleError(w, r, errors.New("foo"))

	if w.Code != 500 {
		t.Fatalf("got: %v", w.Code)
	}

	if act := w.Body.String(); act != "foo\n" {
		t.Fatalf("got: %v", act)
	}
}
