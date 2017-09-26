package main

import (
	"os"

	"github.com/takashabe/go-isucon-exercise/portal/server"
)

func main() {
	cli := &server.CLI{OutStream: os.Stdout, ErrStream: os.Stderr}
	os.Exit(cli.Run(os.Args))
}
