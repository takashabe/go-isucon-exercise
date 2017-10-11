package main

import (
	"os"

	"github.com/takashabe/go-isucon-exercise/bench/benchmark"
)

func main() {
	cli := &benchmark.CLI{OutStream: os.Stdout, ErrStream: os.Stderr}
	os.Exit(cli.Run(os.Args))
}
