package main

// Task implement for each type of benchmark
type Task interface {
	Task(sessions []*Session)
	FinishHook(r Result) Result
}
