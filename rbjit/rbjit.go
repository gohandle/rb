package rbjit

import (
	"errors"
	"net/http"
)

// ErrNoJIT is returned when the response is not JIT
var ErrNoJIT = errors.New("provided response doesn't allow jit callbacks")

// Middleware replaces the default response writer by an implementation that
// allows other components to defer header saving until just before they are written.
// This is usefull for when headers can only be written after a handler has completed, such
// as with session headers. Finally, it will trigger any callbacks also if no
// explicit or implicit write header was performed
type Middleware func(http.Handler) http.Handler

// NewMiddleware creates the actual middleware
func NewMiddleware() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			jitw := New(w).(*jit)

			next.ServeHTTP(jitw, r)

			// if no implit or explit write header was called we trigger
			// the callbacks anyway.
			if !jitw.wroteHeader {
				jitw.trigger()
			}
		})
	}
}

// MustAppendCallback will call AppendCallback but panic when it fails
func MustAppendCallback(w http.ResponseWriter, cb func()) {
	if err := AppendCallback(w, cb); err != nil {
		panic("rbjit: " + err.Error())
	}
}

// AppendCallback will add a callback if 'w' is a jit response
func AppendCallback(w http.ResponseWriter, cb func()) error {
	jw, ok := w.(*jit)
	if !ok {
		return ErrNoJIT
	}

	jw.cbs = append(jw.cbs, cb)
	return nil
}

type jit struct {
	http.ResponseWriter
	cbs         []func()
	wroteHeader bool
}

// New creates a JIT (just-in-time) callback writer
func New(w http.ResponseWriter) http.ResponseWriter {
	return &jit{w, nil, false}
}

// WriteHeader implements the std libray write header but calls
// the just-in-time callback first
func (w *jit) WriteHeader(statusCode int) {
	if w.wroteHeader {
		return
	}

	w.trigger()
	w.ResponseWriter.WriteHeader(statusCode)
	w.wroteHeader = true
}

// Write will call write header if it wasn't called yet
func (w *jit) Write(b []byte) (int, error) {
	if !w.wroteHeader {
		w.WriteHeader(http.StatusOK)
	}

	return w.ResponseWriter.Write(b)
}

func (w *jit) trigger() {
	for _, cb := range w.cbs {
		cb()
	}
}
