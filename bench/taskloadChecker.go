package main

import (
	"math/rand"
	"time"
)

type LoadCheckerTask struct {
	w Worker
}

func (t *LoadCheckerTask) SetWorker(w Worker) {
	t.w = w
}

func (t *LoadCheckerTask) FinishHook(r Result) Result {
	return r
}

func (t *LoadCheckerTask) Task(sessions []*Session) {
	stopAt := time.Now().Add(10000 * time.Millisecond)
	for time.Now().After(stopAt) {
		// TODO implements
		// LoadCheckerTask use 3...9
		rand.Seed(time.Now().UnixNano())
		sub := sessions[3:9]
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
