package main

import (
	"os"

	"github.com/takashabe/go-isucon-exercise/bench/agent"
)

func main() {
	c := agent.CLI{
		OutStream: os.Stdout,
		ErrStream: os.Stderr,
	}
	os.Exit(c.Run(os.Args))
}
