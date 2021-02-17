package rb_test

import (
	"fmt"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"testing"

	rb "github.com/gohandle/rb/rbv2"
	"github.com/gohandle/rb/rbv2/rbgorilla"
	"github.com/gohandle/rb/rbv2/rbjit"
	"github.com/gohandle/rb/rbv2/rbtest"
	"github.com/gorilla/sessions"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest/observer"
)

func TestRequestIDMiddleware(t *testing.T) {
	rnd := rand.New(rand.NewSource(1))
	rb.RandRead = rnd.Read

	t.Run("without any headers", func(t *testing.T) {
		w, r := httptest.NewRecorder(), httptest.NewRequest("GET", "/", nil)
		rb.NewIDMiddleware()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, "rid:%s", rb.RequestID(r.Context()))
		})).ServeHTTP(w, r)

		if w.Body.String() != `rid:Uv38ByGCZU8WP18PmmIdcpVm` {
			t.Fatalf("got: %v", w.Body.String())
		}
	})

	t.Run("without common headers", func(t *testing.T) {
		w, r := httptest.NewRecorder(), httptest.NewRequest("GET", "/", nil)
		r.Header.Set("X-Amzn-Trace-Id", "foo")

		rb.NewIDMiddleware(rb.CommonRequestIDHeaders...)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, "rid:%s", rb.RequestID(r.Context()))
		})).ServeHTTP(w, r)

		if w.Body.String() != `rid:foo` {
			t.Fatalf("got: %v", w.Body.String())
		}
	})
}

func TestRequestLogger(t *testing.T) {
	rnd := rand.New(rand.NewSource(1))
	rb.RandRead = rnd.Read

	lc, obs := observer.New(zap.DebugLevel)

	t.Run("without a request id", func(t *testing.T) {
		w, r := httptest.NewRecorder(), httptest.NewRequest("GET", "/", nil)

		rb.NewLoggerMiddleware(zap.New(lc))(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if act := rb.RequestLogger(r.Context()); act == nil {
				t.Fatalf("got: %v", act)
			}

			rb.L(r).Info("foo")
		})).ServeHTTP(w, r)

		if obs.FilterMessage("foo").Len() != 1 ||
			len(obs.FilterMessage("foo").All()[0].Context) != 5 ||
			obs.FilterMessage("foo").All()[0].Context[0].Key != "request_url" {
			t.Fatalf("got: %v", obs.All())
		}
	})

	t.Run("with a request id", func(t *testing.T) {
		w, r := httptest.NewRecorder(), httptest.NewRequest("GET", "/", nil)

		rb.NewIDMiddleware()(rb.NewLoggerMiddleware(zap.New(lc))(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if act := rb.RequestLogger(r.Context()); act == nil {
				t.Fatalf("got: %v", act)
			}

			rb.L(r).Info("bar")
		}))).ServeHTTP(w, r)

		if obs.FilterMessage("bar").Len() != 1 ||
			obs.FilterMessage("bar").All()[0].Context[5].Key != "request_id" {
			t.Fatalf("got: %v", obs.All())
		}
	})

	// quick test for getting a nop logger
	if act := rb.L(); act == nil {
		t.Fatalf("got:%v", act)
	}
}

func TestSessionSaveMiddleware(t *testing.T) {
	sc := rb.NewSessionCore(rbgorilla.AdaptSessionStore(sessions.NewCookieStore(make([]byte, 32))))

	w, r := httptest.NewRecorder(), httptest.NewRequest("GET", "/", nil)
	rbjit.NewMiddleware()(
		rb.NewSessionSaveMiddleware(sc)(
			http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				sc.Session(w, r).Set("foo", "bar")
			}))).ServeHTTP(w, r)

	s := rbtest.ReadSession(t, sc, rb.DefaultCookieName, w.Header().Get("Set-Cookie"))
	if act := s.Get("foo").(string); act != "bar" {
		t.Fatalf("got: %v", act)
	}
}
