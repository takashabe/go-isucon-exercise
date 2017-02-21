package main

// InitTask is initialize
type InitTask struct {
	w   Worker
	ctx Ctx
}

func (t *InitTask) Task() {
	t.w.getAndCheck(nil, "/initialize", "INITIALIZE", func(c *Checker) {
		c.isStatusCode(200)
		c.isStatusCode(400)
	})
}
