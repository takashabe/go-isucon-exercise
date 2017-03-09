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

func (t *LoadTask) GetWorker() Worker {
	return t.w
}

func (t *LoadTask) FinishHook() Result {
	result := t.w.getResult()
	if len(result.Violations) > 0 {
		result.Valid = false
	}
	return *result
}

func (t *LoadTask) Task(sessions []*Session) {
	timeout := time.After(1000 * time.Millisecond)
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

func (t *LoadTask) run(sessions []*Session) {
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
