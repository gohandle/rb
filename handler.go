package rb

import (
	"net/http"
)

// Handler is close to the standardlib http.Handler but allows for handling errors
type Handler interface {
	ErrServeHTTP(http.ResponseWriter, *http.Request) error
}

// HandlerFunc allows for casting functions to implement our Handler
type HandlerFunc func(http.ResponseWriter, *http.Request) error

// ErrServeHTTP implements our Handler
func (f HandlerFunc) ErrServeHTTP(w http.ResponseWriter, r *http.Request) error {
	return f(w, r)
}

// Action is a short method that turns serves Handlers as a standard library handler but
// uses the apps configured error handlers to render any error that occurs.
func (a *App) Action(f HandlerFunc) http.Handler {
	return a.ActionHandler(f)
}

// ActionHandler is a short method that turns serves Handlers as a standard library handler but
// uses the apps configured error handlers to render any error that occurs.
func (a *App) ActionHandler(next Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := next.ErrServeHTTP(w, r)
		if err == nil {
			return // OK
		}

		a.handleErrorOrPanic(w, r, err)
	})
}
