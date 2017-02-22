package main

// TODO: Taskをtaskフォルダ配下にして、Taskインタフェースを外出しにして他のTaskも書き始める.
// Task implement for each type of benchmark
type Task interface {
	Task()
	FinishHook(r Result) Result
}

// InitTask is initialize
type InitTask struct {
	w Worker
}

func (t *InitTask) Task() {
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
