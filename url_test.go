package rb_test

import (
	"testing"

	"github.com/gohandle/rb"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

func TestURLGeneration(t *testing.T) {
	r := mux.NewRouter()
	r.Name("foo").Path("/foo/{id}/bar")
	r.Name("with_host").Host("localhost:9090").Path("/bar").Schemes("ftp")

	a := rb.New(zap.NewNop(), nil, nil, nil, nil, r)

	if act := a.URL("/foo/bar"); act != "/foo/bar" {
		t.Fatalf("got: %v", act)
	}

	if act := a.URL("foo", rb.URLVar("id", "rab")); act != "/foo/rab/bar" {
		t.Fatalf("got: %v", act)
	}

	if act := a.URL("with_host"); act != "ftp://localhost:9090/bar" {
		t.Fatalf("got: %v", act)
	}

	a.BasePath = "/prod/"
	if act := a.URL("with_host"); act != "ftp://localhost:9090/prod/bar" {
		t.Fatalf("got: %v", act)
	}

}
