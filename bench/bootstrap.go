package main

// BootstrapTask checks initial content consistency
type BootstrapTask struct {
	w Worker
}

func (t *BootstrapTask) Task() {
	t.w.getAndCheck(nil, "/initialize", "INITIALIZE", func(c *Checker) {
		c.isStatusCode(200)
	})
}

func (t *BootstrapTask) FinishHook(r Result) Result {
	r.Valid = true

	if len(r.Violations) > 0 {
		r.Valid = false
	}
	return r
}
