package agent

import (
	"flag"
	"fmt"
	"io"

	"github.com/pkg/errors"
)

// default parameters
const (
	defaultInterval = 5
	defaultPubsub   = "http://localhost:9000"
	defaultHost     = "localhost"
	defaultPort     = 8080
)

// Exit codes. used only in Run()
const (
	ExitCodeOK = 0

	// Specific error codes. begin 10-
	ExitCodeError = 10 + iota
	ExitCodeParseError
	ExitCodeInvalidArgsError
	ExitCodeSetupServerError
)

var (
	// ErrParseFailed is failed to cli args parse
	ErrParseFailed = errors.New("failed to parse args")
)

type param struct {
	interval  int
	pubsub    string
	benchmark string
	param     string
	host      string
	port      int
}

// CLI is the command line interface object
type CLI struct {
	OutStream io.Writer
	ErrStream io.Writer
}

// Run invokes the CLI with the given arguments
func (c *CLI) Run(args []string) int {
	p := &param{}
	err := c.parseArgs(args[1:], p)
	if err != nil {
		fmt.Fprintf(c.ErrStream, "args parse error: %v", err)
		return ExitCodeParseError
	}

	dispather, err := NewDispatch(p.benchmark, p.param, p.host, p.port)
	if err != nil {
		fmt.Fprintf(c.ErrStream, "failed to initialized dispatcher: %v", err)
		return ExitCodeInvalidArgsError
	}
	agent, err := NewAgent(p.interval, p.pubsub, dispather)
	if err != nil {
		fmt.Fprintf(c.ErrStream, "failed to initialized agent: %v", err)
		return ExitCodeInvalidArgsError
	}

	err = agent.Run()
	if err != nil {
		fmt.Fprintf(c.ErrStream, "failed to running agent: %v", err)
		return ExitCodeError
	}
	return ExitCodeOK
}

func (c *CLI) parseArgs(args []string, p *param) error {
	flags := flag.NewFlagSet("param", flag.ContinueOnError)
	flags.SetOutput(c.ErrStream)

	flags.IntVar(&p.interval, "interval", defaultInterval, "Running port. require unused port.")
	flags.StringVar(&p.pubsub, "pubsub", defaultPubsub, "Pubsub Server URL.")
	flags.StringVar(&p.benchmark, "benchmark", "", "Benchmark script path.")
	flags.StringVar(&p.param, "param", "", "Parameter file path.")
	flags.StringVar(&p.host, "host", defaultHost, "Webapp hostname")
	flags.IntVar(&p.port, "port", defaultPort, "Webapp running port")

	err := flags.Parse(args)
	if err != nil {
		return errors.Wrapf(ErrParseFailed, err.Error())
	}
	return nil
}
