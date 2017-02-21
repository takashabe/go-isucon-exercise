package main

import "testing"

func TestTask(t *testing.T) {
	ctx := newCtx()
	worker := Worker{
		task: InitTask{},
	}
}
