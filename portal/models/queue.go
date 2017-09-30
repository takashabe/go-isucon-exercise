package models

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/pkg/errors"
	"github.com/takashabe/go-message-queue/client"
)

// pubsub server names
const (
	pubsubPortal      = "portal"
	pubsubBenchmarker = "result"
)

// Error variables
var (
	ErrExistQueue = errors.New("the queue already exist")
)

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

	portalTopic, err := q.setupTopic(ctx, pubsubPortal)
	if err != nil {
		return nil, err
	}
	_, err = q.setupSubscription(ctx, portalTopic, pubsubPortal)
	if err != nil {
		return nil, err
	}

	benchTopic, err := q.setupTopic(ctx, pubsubBenchmarker)
	if err != nil {
		return nil, err
	}
	_, err = q.setupSubscription(ctx, benchTopic, pubsubBenchmarker)
	if err != nil {
		return nil, err
	}
	return q, nil
}

func (q *Queue) setupTopic(ctx context.Context, id string) (*client.Topic, error) {
	exist, err := q.c.Topic(id).Exists(ctx)
	if err != nil {
		return nil, err
	}
	if exist {
		return nil, nil
	}

	return q.c.CreateTopic(ctx, id)
}

func (q *Queue) setupSubscription(ctx context.Context, topic *client.Topic, id string) (*client.Subscription, error) {
	sub := q.c.Subscription(id)
	exist, err := sub.Exists(ctx)
	if err != nil {
		return nil, err
	}
	if exist {
		return sub, nil
	}

	cfg := client.SubscriptionConfig{Topic: topic}
	return q.c.CreateSubscription(ctx, id, cfg)
}

// Publish send queue message
func (q *Queue) Publish(ctx context.Context, teamID int) (string, error) {
	now := time.Now()
	d, err := NewDatastore()
	if err != nil {
		return "", err
	}
	defer d.Close()

	exist, err := q.isExistQueue(d, teamID)
	if err != nil {
		return "", err
	}
	if exist {
		return "", ErrExistQueue
	}

	result := q.c.Topic(pubsubPortal).Publish(ctx, &client.Message{
		Attributes: map[string]string{"team_id": fmt.Sprintf("%d", teamID)},
	})
	msgID, err := result.Get(ctx)
	if err != nil {
		return "", err
	}
	return msgID, d.saveQueues(teamID, msgID, now)
}

// isExistQueue returns whether finished_at is NULL at the queues
func (q *Queue) isExistQueue(d *Datastore, teamID int) (bool, error) {
	row, err := d.findQueueByTeamID(teamID)
	if err != nil {
		return false, err
	}

	var (
		msgID      string
		finishedAt mysql.NullTime
	)
	err = row.Scan(&msgID, &finishedAt)
	if err == sql.ErrNoRows {
		return false, nil
	}
	return !finishedAt.Valid, err
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
	SourceMessageID string
	Err             error
}

// PullAndSave receive queue message and save message for Datastore
func (q *Queue) PullAndSave(ctx context.Context) error {
	var (
		response QueueResponse
		result   BenchmarkResult
		ackID    string
	)

	sub := q.c.Subscription(pubsubBenchmarker)
	err := sub.Receive(ctx, func(ctx context.Context, msg *client.Message) {
		ackID = msg.AckID

		err := json.NewDecoder(bytes.NewBuffer(msg.Data)).Decode(&result)
		if err != nil {
			response.Err = err
			return
		}
		response.BenchmarkResult = result

		teamID, err := strconv.Atoi(msg.Attributes["team_id"])
		if err != nil {
			response.Err = errors.Wrapf(err, "invalid attributes: team_id")
			return
		}
		response.TeamID = teamID

		unixTime, err := strconv.ParseInt(msg.Attributes["created_at"], 10, 64)
		if err != nil {
			response.Err = errors.Wrapf(err, "invalid attributes: created_at")
			return
		}
		response.CreatedAt = time.Unix(unixTime, 0)

		msgID, ok := msg.Attributes["source_msg_id"]
		if !ok {
			response.Err = errors.New("invalid attributes: source_msg_id")
		}
		response.SourceMessageID = msgID
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

	d, err := NewDatastore()
	if err != nil {
		if nerr := sub.Nack(ctx, []string{ackID}); nerr != nil {
			err = errors.Wrap(err, nerr.Error())
		}
		return err
	}
	defer d.Close()

	err = d.saveScoreAndUpdateQueue(response)
	if err != nil {
		if nerr := sub.Nack(ctx, []string{ackID}); nerr != nil {
			err = errors.Wrap(err, nerr.Error())
		}
		return err
	}

	return sub.Ack(ctx, []string{ackID})
}

// CurrentQueue represent current active queue
type CurrentQueue struct {
	ID     string `json:"message_id"`
	MyTeam bool   `json:"my_team"`
}

// CurrentQueues returns active current queues
func (q *Queue) CurrentQueues(ctx context.Context, teamID int) ([]CurrentQueue, error) {
	d, err := NewDatastore()
	if err != nil {
		return nil, err
	}
	defer d.Close()

	sub := q.c.Subscription(pubsubPortal)
	raw, err := sub.StatsDetail(ctx)
	if err != nil {
		return nil, err
	}

	type jsonMapper struct {
		Messages []string `json:"subscription.portal.current_messages"`
	}
	var decode jsonMapper
	err = json.NewDecoder(bytes.NewBuffer(raw)).Decode(&decode)
	if err != nil {
		return nil, err
	}
	if len(decode.Messages) == 0 {
		return []CurrentQueue{}, nil
	}

	row, err := d.findQueueByTeamID(teamID)
	if err != nil {
		return nil, err
	}
	var (
		msgID      string
		finishedAt mysql.NullTime
	)
	err = row.Scan(&msgID, &finishedAt)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	queues := []CurrentQueue{}
	for _, msg := range decode.Messages {
		queues = append(queues, CurrentQueue{
			ID:     msg,
			MyTeam: msg == msgID,
		})
	}
	return queues, nil
}
