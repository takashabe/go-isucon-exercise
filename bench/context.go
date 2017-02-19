package main

// Context is environment settings for each Worker
type Context struct {
	schema   string
	host     string
	port     int
	agent    string
	sessions []Session
}

// TODO: move to session.go
// Session is save cookies
type Session struct{}
