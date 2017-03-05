package main

import (
	"fmt"
	"net/url"
	"time"

	bench "github.com/takashabe/go-isucon-exercise/bench"
)

func IsuconWorkers() []*bench.Worker {
	init := &InitTask{}
	bootstrap := &BootstrapTask{}
	load := &LoadTask{}
	loadChecker := &LoadCheckerTask{}

	ws := []*Worker{
		NewWorker(),
		// NewWorker().setRunningTime(30000).setTasks(init),
		NewWorker().setRunningTime(60000).setTasks(init, bootstrap),
		NewWorker().setRunningTime(60000).setTasks(load, load, loadChecker),
	}
	return ws
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
