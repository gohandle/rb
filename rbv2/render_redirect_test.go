package rb_test

import (
	"net/http/httptest"
	"testing"

	rb "github.com/gohandle/rb/rbv2"
)

func TestRenderRedirect(t *testing.T) {
	w, r := httptest.NewRecorder(), httptest.NewRequest("GET", "/", nil)
	rb.NewRenderCore(nil).Render(w, r, rb.Redirect("/"))

	if w.Code != 302 {
		t.Fatalf("got: %v", w.Code)
	}

	if w.Header().Get("Location") != "/" {
		t.Fatalf("got: %v", w.Header())
	}
}
