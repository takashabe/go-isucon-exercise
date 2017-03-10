package main

import (
	"flag"
	"io"
	"log"
	"os"

	"github.com/pkg/errors"
)

const (
	defaultHost = "localhost"
	defaultPort = 80
	defaultTime = 3 * 60 * 1000
	defaultFile = "testdata/param.json"
)

// Exit codes. used only in Run()
const (
	ExitCodeOK = 0

	// Specific error codes. begin 10-
	ExitCodeError = 10 + iota
	ExitCodeParseError
)

var (
	// ErrParseFailed is failed to cli args parse
	ErrParseFailed = errors.New("failed to parse args")
)

// PrintDebugf behaves like log.Printf only in the debug env
func PrintDebugf(format string, args ...interface{}) {
	if env := os.Getenv("ISUCON_BENCH_DEBUG"); len(env) != 0 {
		log.Printf("[DEBUG] "+format+"\n", args...)
	}
}

type param struct {
	host string
	port int
	time int
	file string
}

// CLI is the command line interface object
type CLI struct {
	outStream io.Writer
	errStream io.Writer
}

// Run invokes the CLI with the given arguments
func (c *CLI) Run(args []string) int {
	param := &param{}
	c.parseArgs(args, param)
	return ExitCodeOK
}

func (c *CLI) parseArgs(args []string, p *param) error {
	flags := flag.NewFlagSet("bench", flag.ContinueOnError)
	flags.SetOutput(c.errStream)

	flags.StringVar(&p.file, "file", defaultFile, "")
	flags.IntVar(&p.port, "port", defaultPort, "")
	flags.IntVar(&p.time, "time", defaultTime, "")
	flags.StringVar(&p.host, "host", defaultHost, "")

	err := flags.Parse(args)
	if err != nil {
		PrintDebugf("parse error: %v", err)
		return errors.Wrapf(ErrParseFailed, err.Error())
	}
	return nil
}
