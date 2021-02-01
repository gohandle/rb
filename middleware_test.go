package rb

import (
	"fmt"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"testing"

	"go.uber.org/zap"
	"go.uber.org/zap/zaptest/observer"
)

func TestRequestIDMiddleware(t *testing.T) {
	rnd := rand.New(rand.NewSource(1))
	RandRead = rnd.Read

	a := New(zap.NewNop(), nil, nil, nil, nil, nil)

	t.Run("without any headers", func(t *testing.T) {
		w, r := httptest.NewRecorder(), httptest.NewRequest("GET", "/", nil)
		a.IDMiddleware()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, "rid:%s", RequestID(r.Context()))
		})).ServeHTTP(w, r)

		if w.Body.String() != `rid:Uv38ByGCZU8WP18PmmIdcpVm` {
			t.Fatalf("got: %v", w.Body.String())
		}
	})

	t.Run("without common headers", func(t *testing.T) {
		w, r := httptest.NewRecorder(), httptest.NewRequest("GET", "/", nil)
		r.Header.Set("X-Amzn-Trace-Id", "foo")

		a.IDMiddleware(CommonRequestIDHeaders...)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, "rid:%s", RequestID(r.Context()))
		})).ServeHTTP(w, r)

		if w.Body.String() != `rid:foo` {
			t.Fatalf("got: %v", w.Body.String())
		}
	})
}

func TestRequestLogger(t *testing.T) {
	rnd := rand.New(rand.NewSource(1))
	RandRead = rnd.Read

	lc, obs := observer.New(zap.DebugLevel)
	a := New(zap.New(lc), nil, nil, nil, nil, nil)

	t.Run("without a request id", func(t *testing.T) {
		w, r := httptest.NewRecorder(), httptest.NewRequest("GET", "/", nil)

		a.LoggerMiddleware()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if act := RequestLogger(r.Context()); act == nil {
				t.Fatalf("got: %v", act)
			}

			a.L(r).Info("foo")
		})).ServeHTTP(w, r)

		if obs.FilterMessage("foo").Len() != 1 ||
			len(obs.FilterMessage("foo").All()[0].Context) != 5 ||
			obs.FilterMessage("foo").All()[0].Context[0].Key != "request_url" {
			t.Fatalf("got: %v", obs.All())
		}
	})

	t.Run("with a request id", func(t *testing.T) {
		w, r := httptest.NewRecorder(), httptest.NewRequest("GET", "/", nil)

		a.IDMiddleware()(a.LoggerMiddleware()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if act := RequestLogger(r.Context()); act == nil {
				t.Fatalf("got: %v", act)
			}

			a.L(r).Info("bar")
		}))).ServeHTTP(w, r)

		if obs.FilterMessage("bar").Len() != 1 ||
			obs.FilterMessage("bar").All()[0].Context[5].Key != "request_id" {
			t.Fatalf("got: %v", obs.All())
		}
	})

	// quick test for getting a nop logger
	a2 := new(App)
	if act := a2.L(); act == nil {
		t.Fatalf("got:%v", act)
	}
}
