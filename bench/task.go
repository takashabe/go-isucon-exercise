package main

import (
	"fmt"
	"net/url"
	"time"
)

// Task implement for each type of benchmark
type Task interface {
	SetWorker(w Worker)
	GetWorker() Worker
	Task(sessions []*Session)
	FinishHook() Result
}

// Task utilities
var util = taskUtil{}

type taskUtil struct{}

func (t *taskUtil) makeLoginParam(email, password string) url.Values {
	values := url.Values{}
	values.Set("email", email)
	values.Set("password", password)
	return values
}

func (t *taskUtil) makeTweetParam() url.Values {
	p := url.Values{}
	p.Set("content", fmt.Sprint(time.Now()))
	return p
}

func IsuconWorkers() []*Worker {
	init := &InitTask{}
	bootstrap := &BootstrapTask{}
	load := &LoadTask{}
	loadChecker := &LoadCheckerTask{}

	ws := []*Worker{
		NewWorker().setRunningTime(30000).setTasks(init),
		NewWorker().setRunningTime(30000).setTasks(bootstrap),
		NewWorker().setRunningTime(60000).setTasks(load, load, load, loadChecker),
	}
	return ws
}
