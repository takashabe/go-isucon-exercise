package main

// InitTask is initialize
type InitTask struct {
	w Worker
}

func (t *InitTask) SetWorker(w Worker) {
	t.w = w
}

func (t *InitTask) GetWorker() Worker {
	return t.w
}

func (t *InitTask) FinishHook() Result {
	if len(t.w.result.Violations) > 0 {
		t.w.result.Valid = false
	}
	return *t.w.result
}

func (t *InitTask) Task(sessions []*Session) {
	t.w.getAndCheck(nil, "/initialize", "INITIALIZE", func(c *Checker) {
		c.isStatusCode(200)
	})
}
