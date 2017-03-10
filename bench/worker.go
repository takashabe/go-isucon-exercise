package main

import "sync"

// Worker is send requests
// TODO: export check request functions from Worker
type Worker struct {
	ctx   Ctx
	tasks []Task
	mu    sync.Mutex
}

func NewWorker(ctx Ctx, time int, tasks []Task) *Worker {
	ctx.workerRunningTime = time
	return &Worker{
		ctx:   ctx,
		tasks: tasks,
	}
}

func (w *Worker) run() *Result {
	allResult := newResult()
	dones := make(chan Result, len(w.tasks))
	for _, t := range w.tasks {
		go func() {
			driver := &Driver{
				result: newResult(),
				ctx:    w.ctx,
			}
			t.Task(w.ctx, driver)
			r := t.FinishHook(*driver.result)
			dones <- r
		}()
	}
	for i := 0; i < len(w.tasks); i++ {
		r := <-dones
		allResult.Merge(r)
	}
	return allResult
}
