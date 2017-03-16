package main

// InitTask is initialize
type InitTask struct{}

func (t *InitTask) FinishHook(r Result) Result {
	if len(r.Violations) > 0 {
		r.Fail()
	}
	return r
}

func (t *InitTask) Task(ctx Ctx, d *Driver) *Driver {
	d.getAndCheck(nil, "/initialize", "INITIALIZE", func(c *Checker) {
		c.isStatusCode(200)
		c.respondUntil(ctx.workerRunningTime)
	})

	return d
}
