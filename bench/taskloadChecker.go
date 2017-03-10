package main

import (
	"math/rand"
	"time"
)

type LoadCheckerTask struct {
}

func (t *LoadCheckerTask) FinishHook(r Result) Result {
	if len(r.Violations) > 0 {
		r.Fail()
	}
	return r
}

func (t *LoadCheckerTask) Task(ctx Ctx, d *Driver) *Driver {
	timeout := time.After(100 * time.Millisecond)
	// TODO: more interrupt timeout in run(). for each send request
	for {
		select {
		case <-timeout:
			return d
		default:
			t.run(ctx, d)
		}
	}
}

func (t *LoadCheckerTask) run(ctx Ctx, d *Driver) {
	// LoadTask use 10...
	rand.Seed(time.Now().UnixNano())
	sub := ctx.sessions[10:]
	s1 := sub[rand.Intn(len(sub))]
	// s2 := sub[rand.Intn(len(sub))]
	// s3 := sub[rand.Intn(len(sub))]

	d.get(s1, "/logout")
	d.post(s1, "/login", util.makeLoginParam(s1.param.Email, s1.param.Password))
	d.postAndCheck(s1, "/tweet", util.makeTweetParam(), "POST TWEET", func(c *Checker) {
		c.isRedirect("/")
	})
}
