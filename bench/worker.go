package main

import (
	"sync"

	"github.com/k0kubun/pp"
)

// Worker is send requests
// TODO: export check request functions from Worker
type Worker struct {
	ctx    Ctx
	tasks  []Task
	result *Result
	mu     sync.Mutex
}

func NewWorker() *Worker {
	return &Worker{
		ctx:    *newCtx(),
		result: newResult(),
	}
}

func (w *Worker) setRunningTime(t int) *Worker {
	w.ctx.maxRunningTime = t
	return w
}

func (w *Worker) setTasks(tasks ...Task) *Worker {
	w.tasks = tasks
	return w
}

func (w *Worker) getResult() *Result {
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.result
}

func (w *Worker) run(sessions []*Session) *Result {
	result := newResult()
	dones := make(chan struct{}, len(w.tasks))
	for _, t := range w.tasks {
		go func() {
			t.SetWorker(*w)
			t.Task(sessions)
			r := t.FinishHook()
			pp.Println(r)
			result.Merge(r)
			dones <- struct{}{}
		}()
	}
	for i := 0; i < len(w.tasks); i++ {
		<-dones
	}
	return result
}
