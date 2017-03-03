package main

import (
	"lib/bench"
	"os"
)

func main() {
	os.Exit(bench.Run(os.Args))
}
