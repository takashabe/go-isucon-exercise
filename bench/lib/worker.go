package bench

// Worker is send requests
// TODO: export check request functions from Worker
type Worker struct {
	ctx    Ctx
	tasks  []Task
	result *Result
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
