package rb

import (
	"net/http"

	"github.com/gorilla/sessions"
)

type Session interface {
	Set(k, v interface{}) Session
	Del(k interface{}) Session
	Get(k interface{}) (v interface{})
}

type session struct {
	s *sessions.Session
	r *http.Request
	w http.ResponseWriter
}

func (s session) save() {
	if err := s.s.Save(s.r, s.w); err != nil {
		// @TODO log any session save error
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

	s, _ := a.sess.Get(r, o.sessionName)
	// @TODO log session getting errors

	return session{s, r, w}
}
