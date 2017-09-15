package models

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/pkg/errors"
	"github.com/takashabe/go-message-queue/client"
)

// PubsubServerName represent pubsub topic and subscription name
const PubsubServerName = "benchmark"

// Queue is client for the pubsub
type Queue struct {
	c *client.Client
}

// NewQueue returns initialized Queue
func NewQueue(addr string) (*Queue, error) {
	ctx := context.Background()
	client, err := client.NewClient(ctx, addr)
	if err != nil {
		return nil, err
	}
	q := &Queue{
		c: client,
	}

	topic, err := q.setupTopic(ctx)
	if err != nil {
		return nil, err
	}
	_, err = q.setupSubscription(ctx, topic)
	if err != nil {
		return nil, err
	}
	return q, nil
}

func (q *Queue) setupTopic(ctx context.Context) (*client.Topic, error) {
	exist, err := q.c.Topic(PubsubServerName).Exists(ctx)
	if err != nil {
		return nil, err
	}
	if exist {
		return nil, nil
	}

	return q.c.CreateTopic(ctx, PubsubServerName)
}

func (q *Queue) setupSubscription(ctx context.Context, topic *client.Topic) (*client.Subscription, error) {
	sub := q.c.Subscription(PubsubServerName)
	exist, err := sub.Exists(ctx)
	if err != nil {
		return nil, err
	}
	if exist {
		return sub, nil
	}

	cfg := client.SubscriptionConfig{Topic: topic}
	return q.c.CreateSubscription(ctx, PubsubServerName, cfg)
}

// Publish send queue message
func (q *Queue) Publish(ctx context.Context, teamID int) error {
	result := q.c.Topic(PubsubServerName).Publish(ctx, &client.Message{
		Attributes: map[string]string{"team_id": fmt.Sprintf("%d", teamID)},
	})
	_, err := result.Get(ctx)
	return err
}

// BenchmarkResult represent the benchmark result JSON
type BenchmarkResult struct {
	Valid        bool `json:"valid"`
	RequestCount int  `json:"request_count"`
	ElapsedTime  int  `json:"elapsed_time"`
	Response     struct {
		Success     int `json:"success"`
		Redirect    int `json:"redirect"`
		ClientError int `json:"client_error"`
		ServerError int `json:"server_error"`
		Exception   int `json:"exception"`
	} `json:"response"`
	Violations []struct {
		RequestName string `json:"request_type"`
		Cause       string `json:"description"`
		Count       int    `json:"num"`
	} `json:"violations"`
}

// QueueResponse represent the message of the whole receive queue
type QueueResponse struct {
	TeamID          int
	BenchmarkResult BenchmarkResult
	CreatedAt       time.Time
	Err             error
}

// PullAndSave receive queue message and save message for Datastore
func (q *Queue) PullAndSave(ctx context.Context) error {
	var (
		response QueueResponse
		result   BenchmarkResult
		ackID    string
	)

	sub := q.c.Subscription(PubsubServerName)
	err := sub.Receive(ctx, func(ctx context.Context, msg *client.Message) {
		ackID = msg.AckID

		err := json.NewDecoder(bytes.NewBuffer(msg.Data)).Decode(&result)
		if err != nil {
			response.Err = err
			return
		}
		response.BenchmarkResult = result

		id, err := strconv.Atoi(msg.Attributes["team_id"])
		if err != nil {
			response.Err = err
			return
		}
		response.TeamID = id

		unixTime, err := strconv.ParseInt(msg.Attributes["created_at"], 10, 64)
		if err != nil {
			response.Err = err
			return
		}
		response.CreatedAt = time.Unix(unixTime, 0)
	})
	if err != nil {
		return err
	}

	if response.Err != nil {
		err := response.Err
		if nerr := sub.Nack(ctx, []string{ackID}); nerr != nil {
			err = errors.Wrap(err, nerr.Error())
		}
		return err
	}

	d, err := newDatastore()
	if err != nil {
		if nerr := sub.Nack(ctx, []string{ackID}); nerr != nil {
			err = errors.Wrap(err, nerr.Error())
		}
		return err
	}
	err = d.saveScore(response)
	if err != nil {
		if nerr := sub.Nack(ctx, []string{ackID}); nerr != nil {
			err = errors.Wrap(err, nerr.Error())
		}
		return err
	}

	return sub.Ack(ctx, []string{ackID})
}
