package rbgorilla

import (
	"net/http"

	rb "github.com/gohandle/rb/rbv2"
	"github.com/gorilla/sessions"
	"go.uber.org/zap"
)

type store struct{ sessions.Store }

// AdaptSessionStore will make a gorilla session store available for session access
// in the rb app
func AdaptSessionStore(ss sessions.Store) rb.SessionStore {
	return store{ss}
}

// func (ss store) DecodeSession(name, value string, dst interface{}) error {
// 	sess := sessions.NewSession(ss.Store, name)

// 	// if err := s.Codecs[0].Decode(name, c.Value, &sess.Values); err != nil {
// 	// 	fatalf(tb, "failed to decode cookie: %v", err)
// 	// }
// }

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

func (s session) save() {
	if err := s.s.Save(s.r, s.w); err != nil {
		rb.L(s.r).Error("failed to save session", zap.Error(err))
	}
}

func (s session) Del(k interface{}) rb.Session {
	defer s.save()
	delete(s.s.Values, k)
	return s
}

func (s session) Set(k, v interface{}) rb.Session {
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
