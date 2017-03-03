package bench

import (
	"net/http"
	"net/http/cookiejar"
)

// Session is save cookies
type Session struct {
	cookie http.CookieJar
	param  UserSchema
}

func newSession(p UserSchema) *Session {
	c, err := cookiejar.New(nil)
	if err != nil {
		return nil
	}
	return &Session{
		cookie: c,
		param:  p,
	}
}
