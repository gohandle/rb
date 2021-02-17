package rbgorilla

import (
	"fmt"
	"net/http"

	rb "github.com/gohandle/rb"
	"github.com/gorilla/sessions"
)

type store struct{ sessions.Store }

// AdaptSessionStore will make a gorilla session store available for session access
// in the rb app
func AdaptSessionStore(ss sessions.Store) rb.SessionStore {
	return store{ss}
}

func (ss store) SaveSession(w http.ResponseWriter, r *http.Request, s rb.Session) error {
	sess, ok := s.(session)
	if !ok {
		return fmt.Errorf("save session received a session that was not created in the same store")
	}

	return ss.Store.Save(r, w, sess.s)
}

func (ss store) LoadSession(w http.ResponseWriter, r *http.Request, name string) (rb.Session, error) {
	s, err := ss.Store.Get(r, name)
	if err != nil {
		return nil, err
	}

	return session{s, w, r}, nil
}

type session struct {
	s *sessions.Session
	w http.ResponseWriter
	r *http.Request
}

func (s session) Del(k interface{}) rb.Session {
	delete(s.s.Values, k)
	return s
}

func (s session) Set(k, v interface{}) rb.Session {
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
