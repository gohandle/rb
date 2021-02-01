package rb

import (
	"net/http"

	"github.com/gorilla/sessions"
	"go.uber.org/zap"
)

type SessionReader interface {
	Get(k interface{}) (v interface{})
}

type Session interface {
	Set(k, v interface{}) Session
	Del(k interface{}) Session
	Pop(k interface{}) (v interface{})
	SessionReader
}

type session struct {
	s *sessions.Session
	r *http.Request
	w http.ResponseWriter
	a *App
}

func (s session) save() {
	if err := s.s.Save(s.r, s.w); err != nil {
		s.a.L(s.r).Error("failed to save session", zap.Error(err))
	}
}

func (s session) Del(k interface{}) Session {
	defer s.save()
	delete(s.s.Values, k)
	return s
}

func (s session) Set(k, v interface{}) Session {
	defer s.save()
	s.s.Values[k] = v
	return s
}

func (s session) Get(k interface{}) (v interface{}) {
	v, _ = s.s.Values[k]
	return
}

func (s session) Pop(k interface{}) interface{} {
	v, ok := s.s.Values[k]
	if ok {
		s.Del(k)
	}

	return v
}

type sessionOpts struct {
	sessionName string
}

type SessionOption func(*sessionOpts)

func SessionName(n string) SessionOption {
	return func(o *sessionOpts) {
		o.sessionName = n
	}
}

var DefaultSessionName = "rb"

func (a *App) Session(w http.ResponseWriter, r *http.Request, opts ...SessionOption) Session {
	var o sessionOpts
	for _, opt := range opts {
		opt(&o)
	}

	if o.sessionName == "" {
		o.sessionName = DefaultSessionName
	}

	s, err := a.sess.Get(r, o.sessionName)
	if err != nil {
		a.L(r).Error("failed to read session, continue with new one", zap.Error(err))
	}

	return session{s, r, w, a}
}
