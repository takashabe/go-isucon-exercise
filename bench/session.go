package main

import "net/http"

// Session is save cookies
type Session struct {
	cookie http.CookieJar
}
