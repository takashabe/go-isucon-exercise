package main

import (
	"math/rand"
	"time"
)

type LoadCheckerTask struct {
	timeout time.Time
}

func (t *LoadCheckerTask) FinishHook(r Result) Result {
	if len(r.Violations) > 0 {
		r.Fail()
	}
	return r
}

func (t *LoadCheckerTask) Task(ctx Ctx, d *Driver) *Driver {
	runningTime := ctx.workerRunningTime
	t.timeout = time.Now().Add(time.Millisecond * time.Duration(runningTime))
	for {
		if t.isTimeout() {
			return d
		}
		t.run(ctx, d)
	}
}

func (t *LoadCheckerTask) isTimeout() bool {
	return t.timeout.Before(time.Now())
}

func (t *LoadCheckerTask) run(ctx Ctx, d *Driver) {
	// LoadTask use 10...
	rand.Seed(time.Now().UnixNano())
	sub := ctx.sessions[10:]
	s1 := sub[rand.Intn(len(sub))]
	s2 := sub[rand.Intn(len(sub))]
	s3 := sub[rand.Intn(len(sub))]

	s1.lockFunc(func() {
		d.get(s1, "/logout")
		d.post(s1, "/login", util.makeLoginParam(s1.param.Email, s1.param.Password))
		d.postAndCheck(s1, "/tweet", util.makeTweetParam(), "POST TWEET", func(c *Checker) {
			c.isRedirect("/")
			// TODO: responseUntil(ctx.postTimeout)
		})
	})
	if t.isTimeout() {
		return
	}

	s2.lockFunc(func() {
		d.get(s2, "/logout")
		d.post(s2, "/login", util.makeLoginParam(s2.param.Email, s2.param.Password))
		d.getAndCheck(s2, "/following", "GET FOLLOWING", func(c *Checker) {
			c.isStatusCode(200)
		})
	})
	if t.isTimeout() {
		return
	}

	s3.lockFunc(func() {
		d.get(s3, "/logout")
		d.post(s3, "/login", util.makeLoginParam(s3.param.Email, s3.param.Password))
		d.getAndCheck(s3, "/followers", "GET FOLLOWERS", func(c *Checker) {
			c.isStatusCode(200)
		})
	})
	if t.isTimeout() {
		return
	}
}
