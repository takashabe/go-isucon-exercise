package main

import (
	"fmt"
	"net/url"
	"time"
)

// Task implement for each type of benchmark
type Task interface {
	Task(ctx Ctx, driver *Driver) *Driver
	FinishHook(result Result) Result
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

// Specific workers
func IsuconWorkOrder() []*WorkOrder {
	init := &InitTask{}
	bootstrap := &BootstrapTask{}
	loadChecker := &LoadCheckerTask{}
	load := &LoadTask{}

	order := []*WorkOrder{
		{30000, []Task{init}},
		{30000, []Task{bootstrap}},
		{1000, []Task{load, load, load, load, load, loadChecker}},
	}
	return order
}
