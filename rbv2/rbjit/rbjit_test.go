package rbjit

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestJIT(t *testing.T) {
	t.Run("implicit write", func(t *testing.T) {
		w, r := httptest.NewRecorder(), httptest.NewRequest("GET", "/", nil)
		var calls string
		NewMiddleware()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			AppendCallback(w, func() { calls += "a" })
			AppendCallback(w, func() { calls += "b" })
			fmt.Fprintf(w, "foo")
			MustAppendCallback(w, func() { calls += "c" }) // this shouldn't do anything
			fmt.Fprintf(w, "bar")                          // no callbacks
		})).ServeHTTP(w, r)

		if calls != "ab" {
			t.Fatalf("got: %v", calls)
		}
	})

	t.Run("explicit write", func(t *testing.T) {
		w, r := httptest.NewRecorder(), httptest.NewRequest("GET", "/", nil)
		var calls string
		NewMiddleware()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			AppendCallback(w, func() { calls += "a" })
			AppendCallback(w, func() { calls += "b" })
			w.WriteHeader(201)
			MustAppendCallback(w, func() { calls += "c" }) // this shouldn't do anything
			w.WriteHeader(302)                             // no callbacks
		})).ServeHTTP(w, r)

		if w.Code != 201 {
			t.Fatalf("got: %v", w.Code)
		}

		if act := w.Body.String(); act != "" {
			t.Fatalf("got: %v", act)
		}

		if calls != "ab" {
			t.Fatalf("got: %v", calls)
		}
	})
}
