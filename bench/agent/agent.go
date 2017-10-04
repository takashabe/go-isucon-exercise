package agent

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"

	"github.com/pkg/errors"
	"github.com/takashabe/go-message-queue/client"
)

const (
	pullServerName    = "portal"
	publishServerName = "result"
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
		pullServer:    pullServerName,
		publishServer: publishServerName,
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
func (a *Agent) Run() {
	for {
		ctx := context.Background()
		msg, err := a.Polling(ctx)
		if err != nil {
			log.Println(err.Error())
			continue
		}
		err = a.pubsub.Subscription(a.pullServer).Ack(ctx, []string{msg.AckID})
		if err != nil {
			log.Println(err.Error())
		}

		data, err := a.Dispatch(ctx)
		if err != nil {
			log.Println(err.Error())
		}

		err = a.SendResult(ctx, data, map[string]string{
			"source_msg_id": msg.ID,
			"team_id":       msg.Attributes["team_id"],
			"created_at":    fmt.Sprintf("%d", time.Now().Unix()),
		})
		if err != nil {
			log.Println(err.Error())
		}
	}
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
		var (
			rm       receiveMessage
			isFinish bool
		)
		for {
			err := sub.Receive(ctx, func(ctx context.Context, msg *client.Message) {
				rm.message = msg
				ch <- rm
				isFinish = true
			})
			if err != nil && err != client.ErrNotFoundMessage {
				rm.err = err
				ch <- rm
				isFinish = true
			}
			if isFinish {
				return
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
		opts := make([]string, 0)
		opts = append(opts, "-host="+a.dispatch.host)
		opts = append(opts, "-file="+a.dispatch.paramFile)
		opts = append(opts, fmt.Sprintf("-port=%d", a.dispatch.port))
		cmd := exec.Command(a.dispatch.script, opts...)
		cmd.Stderr = os.Stderr
		cmd.Stdout = os.Stdout
		res, err := cmd.Output()
		ch <- dispatchResponse{data: res, err: err}
	}(drCh)

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case dr := <-drCh:
		return dr.data, dr.err
	}
}

// SendResult benchmark result send to pubsub server
func (a *Agent) SendResult(ctx context.Context, data []byte, attr map[string]string) error {
	res := a.pubsub.Topic(a.publishServer).Publish(ctx, &client.Message{
		Data:       data,
		Attributes: attr,
	})
	_, err := res.Get(ctx)
	return err
}
