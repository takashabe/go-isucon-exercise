package main

type Checker struct {
	worker Worker
	ctx    Context
	result *Result
}

func newChecker(w Worker, c Context, r Result) *Checker {
	return &Checker{
		worker: w,
		ctx:    c,
		result: r,
	}
}

func (c *Checker) isStatusCode(code int) {
}
