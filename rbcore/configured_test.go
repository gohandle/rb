package rbcore_test

import (
	"testing"

	"github.com/gohandle/rb"
	"github.com/gohandle/rb/rbcore"
	"github.com/gohandle/rb/rbgorilla"
	"github.com/gorilla/mux"
)

func TestConfiguredCore(t *testing.T) {
	rc := rbgorilla.AdaptRouter(mux.NewRouter())
	c := rbcore.Configured(rbcore.New(rc, nil, nil, nil, nil, nil), rbcore.Options{
		URLOptions: []rb.URLOption{rb.BasePath("/base")},
	})

	if act := c.URL(nil, nil, "/foo"); act != "/base/foo" {
		t.Fatalf("got: %v", act)
	}

	if act := c.URL(nil, nil, "/foo", rb.BasePath("/foo")); act != "/foo/foo" {
		t.Fatalf("got: %v", act)
	}
}
