package main

type Worker struct {
	task Task
	ctx  Context
}

// Need subclass
type Task interface {
}
