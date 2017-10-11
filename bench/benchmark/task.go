package benchmark

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
	order := []*WorkOrder{
		{30 * time.Second, []Task{&InitTask{}}},
		{30 * time.Second, []Task{&BootstrapTask{}}},
		{2 * time.Second, []Task{&LoadTask{}, &LoadTask{}, &LoadCheckerTask{}}},
	}
	return order
}
