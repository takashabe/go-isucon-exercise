package main

import (
	"fmt"
	"reflect"
	"sync"
)

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

func (w *Worker) String() string {
	s := ""
	s = fmt.Sprintf("Runningtime: %d:\n", w.ctx.workerRunningTime)
	for _, t := range w.tasks {
		s = s + fmt.Sprintf("\t%s", reflect.TypeOf(t).Elem().Name())
	}
	return s
}
