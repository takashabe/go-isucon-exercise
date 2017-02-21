package main

import "testing"

func TestNewCtx(t *testing.T) {
	a := newCtx()
	a.agent = "test"

	b := newCtx()

	if a != b {
		t.Errorf("a=%v, b=%v", a, b)
	}
}
