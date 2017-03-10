package main

import (
	"net/http"
	"net/http/cookiejar"
)

// Session is save cookies
type Session struct {
	cookie http.CookieJar
	param  UserSchema
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
