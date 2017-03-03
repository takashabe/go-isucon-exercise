package task

// InitTask is initialize
type InitTask struct {
	w Worker
}

func (t *InitTask) Task(sessions []*Session) {
	t.w.getAndCheck(nil, "/initialize", "INITIALIZE", func(c *Checker) {
		c.isStatusCode(200)
	})
}

func (t *InitTask) FinishHook(r Result) Result {
	r.Valid = true

	if len(r.Violations) > 0 {
		r.Valid = false
	}
	return r
}
