package rb

import (
	"net/http"

	"go.uber.org/zap"
)

// ActionFunc implements an application action
type ActionFunc func(Ctx) error

// App provides application wide functionality
type App interface {
	Action(af ActionFunc) http.Handler
}

// L is utility method that returns a zap logger, if possible from the request
// context. If not it will return a NopLogger
func L(r ...*http.Request) *zap.Logger {
	if len(r) < 1 {
		return zap.NewNop()
	}

	l := RequestLogger(r[0].Context())
	if l != nil {
		return l
	}

	l, _ = zap.NewDevelopment()
	return l
}
