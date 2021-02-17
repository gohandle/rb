package rb

import (
	"net/http"

	"go.uber.org/zap"
)

// SessionReader is implemented by sessions that are read-only
type SessionReader interface {
	Get(k interface{}) (v interface{})
}

// Session interface is implemented by sessions that can be written to
type Session interface {
	Set(k, v interface{}) Session
	Del(k interface{}) Session
	Pop(k interface{}) (v interface{})
	SessionReader
}

// SessionOpts will hold all options for session control
type SessionOpts struct{}

// SessionOption allows for configuring session handling
type SessionOption func(*SessionOpts)

type sessionCore struct{ SessionStore }

// NewSessionCore creates the session core part of the core
func NewSessionCore(sess SessionStore) SessionCore {
	return &sessionCore{sess}
}

// DefaultSessionName defines how cookies are named for rb applications by default
var DefaultSessionName = "rb"

func (sc *sessionCore) Session(w http.ResponseWriter, r *http.Request, opts ...SessionOption) Session {
	var o SessionOpts
	for _, opt := range opts {
		opt(&o)
	}

	s, err := sc.SessionStore.LoadSession(w, r, DefaultSessionName)
	if err != nil {
		L(r).Error("failed to load session", zap.Error(err), zap.String("cookie_name", DefaultSessionName))
	}

	return s
}
