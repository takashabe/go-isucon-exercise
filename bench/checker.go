package main

import "net/http"

type Checker struct {
	ctx         Ctx
	result      *Result
	requestName string
	response    http.Response
}

func newChecker(c Ctx, r *Result) *Checker {
	return &Checker{
		ctx:    c,
		result: r,
	}
}

func (c *Checker) isStatusCode(code int) {
}
