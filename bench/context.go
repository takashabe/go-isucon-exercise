package main

import "fmt"

// Ctx is environment settings in each Worker
type Ctx struct {
	// url query parameters
	schema string
	host   string
	port   int
	agent  string

	// http client parameters
	getTimeout  int
	postTimeout int

	// benchmark limitation
	maxRunningTime int

	sessions []Session
}

var defaultCtx = Ctx{
	schema:         "http",
	host:           "localhost",
	port:           80,
	agent:          "isucon",
	getTimeout:     30 * 1000,
	postTimeout:    30 * 1000,
	maxRunningTime: 3 * 60 * 1000,
}

func newCtx() *Ctx {
	ctx := defaultCtx
	return &ctx
}

func (c *Ctx) uri(path string) string {
	return fmt.Sprintf("%s://%s:%d%s", c.schema, c.host, c.port, path)
}
