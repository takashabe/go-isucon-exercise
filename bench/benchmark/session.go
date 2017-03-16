package main

import (
	"net/http"
	"net/http/cookiejar"
	"sync"
)

// Session is save cookies
type Session struct {
	cookie http.CookieJar
	param  UserSchema
	mu     sync.Mutex
}

func newSession(p UserSchema) (*Session, error) {
	c, err := cookiejar.New(nil)
	if err != nil {
		return nil, err
	}
	return &Session{
		cookie: c,
		param:  p,
	}, nil
}

func (s *Session) lockFunc(f func()) {
	s.mu.Lock()
	defer s.mu.Unlock()
	f()
}
