package server

import (
	"flag"
	"fmt"
	"io"
)

const (
	defaultPubsubAddr = "localhost:8080"
	defaultPort       = 8080
)

// Exit codes. used only in Run()
const (
	ExitCodeOK = 0

	// Specific error codes. begin 10-
	ExitCodeError = 10 + iota
	ExitCodeParseError
	ExitCodeInvalidArgsError
)

type param struct {
	pubsubAddr string
	port       int
}

// CLI is the command line interface object
type CLI struct {
	OutStream io.Writer
	ErrStream io.Writer
}

// Run invokes the CLI with the given arguments
func (c *CLI) Run(args []string) int {
	param := &param{}
	err := c.parseArgs(args[1:], param)
	if err != nil {
		fmt.Fprintf(c.ErrStream, "args parse error: %#v", err)
		return ExitCodeParseError
	}

	server, err := NewServer(param.pubsubAddr)
	if err != nil {
		fmt.Fprintf(c.ErrStream, "invalid args. failed to initialize server: %#v", err)
		return ExitCodeInvalidArgsError
	}

	err = server.Run(param.port)
	if err != nil {
		fmt.Fprintf(c.ErrStream, "failed to server run: %#v", err)
		return ExitCodeError
	}
	return ExitCodeOK
}

func (c *CLI) parseArgs(args []string, p *param) error {
	flags := flag.NewFlagSet("portal", flag.ContinueOnError)
	flags.SetOutput(c.ErrStream)

	flags.IntVar(&p.port, "port", defaultPort, "")
	flags.StringVar(&p.pubsubAddr, "pubsub", defaultPubsubAddr, "")

	return flags.Parse(args)
}
