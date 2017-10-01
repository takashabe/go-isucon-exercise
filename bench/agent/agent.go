package agent

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/pkg/errors"
	"github.com/takashabe/go-message-queue/client"
)

const (
	defaultInterval   = 5
	pullServerName    = "portal"
	publishServerName = "benchmark"
)

// Agent represent agent configuration
type Agent struct {
	interval      time.Duration
	pubsub        *client.Client
	dispatch      *Dispatch
	pullServer    string
	publishServer string
}

// Dispatch represent benchmark script configuration
type Dispatch struct {
	script    string
	paramFile string
	host      string
	port      int
}

// NewAgent returns initialized Agent
func NewAgent(interval int, pubsub string, dispatch *Dispatch) (*Agent, error) {
	if interval <= 0 {
		interval = defaultInterval
	}

	ctx := context.Background()
	client, err := client.NewClient(ctx, pubsub)
	if err != nil {
		return nil, err
	}
	_, err = client.Stats(ctx)
	if err != nil {
		return nil, err
	}

	return &Agent{
		interval:      time.Duration(interval) * time.Second,
		pubsub:        client,
		dispatch:      dispatch,
		pullServer:    "portal",
		publishServer: "benchmark",
	}, nil
}

// NewDispatch returns initialized Dispatch
func NewDispatch(script, param, host string, port int) (*Dispatch, error) {
	if !isExistFile(script) {
		return nil, errors.New("Not found script file path")
	}
	if !isExistFile(param) {
		return nil, errors.New("Not found param file path")
	}

	return &Dispatch{
		script:    script,
		paramFile: param,
		host:      host,
		port:      port,
	}, nil
}

func isExistFile(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		return os.IsExist(err)
	}
	return true
}

// Run exec polling and dispatch queues
func (a *Agent) Run() error {
	ctx := context.Background()
	msg, err := a.Polling(ctx)
	if err != nil {
		return err
	}
	pp.Println(msg)

	return nil
}

// Polling trying pull until receive the message
func (a *Agent) Polling(ctx context.Context) (*client.Message, error) {
	type receiveMessage struct {
		message *client.Message
		err     error
	}

	rmCh := make(chan receiveMessage)
	go func(ch chan receiveMessage) {
		sub := a.pubsub.Subscription(a.pullServer)
		var rm receiveMessage
		for {
			err := sub.Receive(ctx, func(ctx context.Context, msg *client.Message) {
				rm.message = msg
				ch <- rm
			})
			if err != nil && err != client.ErrNotFoundMessage {
				rm.err = err
				ch <- rm
			}
			time.Sleep(a.interval)
		}
	}(rmCh)

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case rm := <-rmCh:
		return rm.message, rm.err
	}
}

// Dispatch dispatch benchmark request.
// Returns result from the benchmark script.
func (a *Agent) Dispatch(ctx context.Context) ([]byte, error) {
	type dispatchResponse struct {
		data []byte
		err  error
	}

	drCh := make(chan dispatchResponse)
	go func(ch chan dispatchResponse) {
		opt := fmt.Sprintf("-host=%s -file=%s", a.dispatch.host, a.dispatch.paramFile)
		res, err := exec.Command(a.dispatch.script, opt).Output()
		ch <- dispatchResponse{data: res, err: err}
	}(drCh)

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case dr := <-drCh:
		return dr.data, dr.err
	}
}
