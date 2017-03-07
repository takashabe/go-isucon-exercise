package main

import (
	"math/rand"
	"time"
)

type LoadTask struct {
	w Worker
}

func (t *LoadTask) SetWorker(w Worker) {
	t.w = w
}

func (t *LoadTask) FinishHook(r Result) Result {
	r.Valid = true

	if len(r.Violations) > 0 {
		r.Valid = false
	}
	return r
}

func (t *LoadTask) Task(sessions []*Session) {
	stopAt := time.Now().Add(10000 * time.Millisecond)
	for time.Now().After(stopAt) {
		// TODO implements
		// LoadTask use 10...
		rand.Seed(time.Now().UnixNano())
		sub := sessions[10:]
		s1 := sub[rand.Intn(len(sub))]
		// s2 := sub[rand.Intn(len(sub))]
		// s3 := sub[rand.Intn(len(sub))]

		t.w.get(s1, "/logout")
		t.w.post(s1, "/login", util.makeLoginParam(s1.param.Email, s1.param.Password))
		t.w.postAndCheck(s1, "/tweet", util.makeTweetParam(), "POST TWEET", func(c *Checker) {
			c.isRedirect("/")
		})
	}
}
