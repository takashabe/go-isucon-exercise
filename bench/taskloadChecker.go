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

func (t *LoadCheckerTask) GetWorker() Worker {
	return t.w
}

func (t *LoadCheckerTask) FinishHook() Result {
	if len(t.w.result.Violations) > 0 {
		t.w.result.Valid = false
	}
	return *t.w.result
}

func (t *LoadCheckerTask) Task(sessions []*Session) {
	timeout := time.After(100 * time.Millisecond)
	// TODO: more interrupt timeout in run(). for each send request
	for {
		select {
		case <-timeout:
			return
		default:
			t.run(sessions)
		}
	}
}

func (t *LoadCheckerTask) run(sessions []*Session) {
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
