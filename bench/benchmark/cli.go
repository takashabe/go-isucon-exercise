package benchmark

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/pkg/errors"
)

const (
	defaultHost  = "localhost"
	defaultPort  = 80
	defaultFile  = "data/param.json"
	defaultAgent = "isucon_go"
)

// Exit codes. used only in Run()
const (
	ExitCodeOK = 0

	// Specific error codes. begin 10-
	ExitCodeError = 10 + iota
	ExitCodeParseError
	ExitCodeInvalidArgsError
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
	host  string
	port  int
	time  int
	file  string
	agent string
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

	master, err := NewMaster(param.host, param.port, param.file, param.agent)
	if err != nil {
		fmt.Fprintf(c.ErrStream, "invalid args. failed to initialize master: %#v", err)
		return ExitCodeInvalidArgsError
	}

	result, err := master.start()
	if err != nil {
		fmt.Fprintf(c.ErrStream, "failed to benchmark run: %#v", err)
		return ExitCodeError
	}
	fmt.Fprintln(c.OutStream, result)
	return ExitCodeOK
}

func (c *CLI) parseArgs(args []string, p *param) error {
	flags := flag.NewFlagSet("bench", flag.ContinueOnError)
	flags.SetOutput(c.ErrStream)

	flags.StringVar(&p.file, "file", defaultFile, "")
	flags.IntVar(&p.port, "port", defaultPort, "")
	flags.StringVar(&p.host, "host", defaultHost, "")
	flags.StringVar(&p.agent, "agent", defaultAgent, "")

	err := flags.Parse(args)
	if err != nil {
		return errors.Wrapf(ErrParseFailed, err.Error())
	}
	return nil
}
