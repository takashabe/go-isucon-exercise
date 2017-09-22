package models

import (
	"context"
	"fmt"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	"github.com/takashabe/go-message-queue/client"
)

func publishDummyBenchmarkResult(t *testing.T, ts *httptest.Server, teamID int, msgID string, payload []byte) {
	q, err := NewQueue(ts.URL)
	if err != nil {
		t.Fatalf("want non error, got %v", err)
	}
	ctx := context.Background()
	res := q.c.Topic(pubsubBenchmarker).Publish(ctx, &client.Message{
		Data: payload,
		Attributes: map[string]string{
			"team_id":       fmt.Sprintf("%d", teamID),
			"created_at":    fmt.Sprintf("%d", time.Now().Unix()),
			"source_msg_id": msgID,
		},
	})
	_, err = res.Get(ctx)
	if err != nil {
		t.Fatalf("want non error, got %v", err)
	}
}

func publishDummyBenchmarkQueue(t *testing.T, ts *httptest.Server, teamID int) string {
	q, err := NewQueue(ts.URL)
	if err != nil {
		t.Fatalf("want non error, got %v", err)
	}
	ctx := context.Background()
	msgID, err := q.Publish(ctx, teamID)
	if err != nil {
		t.Fatalf("want non error, got %v", err)
	}
	return msgID
}

func TestNewQueue(t *testing.T) {
	ts := setupPubsub(t)
	defer ts.Close()

	// check exist / non exist pubsub components
	for i := 0; i < 2; i++ {
		_, err := NewQueue(ts.URL)
		if err != nil {
			t.Fatalf("want non error, got %v", err)
		}
	}
}

func TestPublish(t *testing.T) {
	setupFixture(t, "fixture/teams.yaml", "fixture/queues.yaml")
	ts := setupPubsub(t)
	defer ts.Close()

	inputTeam := 1
	q, err := NewQueue(ts.URL)
	if err != nil {
		t.Fatalf("want non error, got %v", err)
	}
	ctx := context.Background()
	_, err = q.Publish(ctx, inputTeam)
	if err != nil {
		t.Fatalf("want non error, got %v", err)
	}
	d, err := NewDatastore()
	if err != nil {
		t.Fatalf("want non error, got %v", err)
	}
	_, err = d.queryRow("select * from queues where team_id=?", inputTeam)
	if err != nil {
		t.Fatalf("want non error, got %v", err)
	}

	_, err = q.Publish(ctx, inputTeam)
	if err != ErrExistQueue {
		t.Fatalf("want error %v, got %v", ErrExistQueue, err)
	}
}

func TestPull(t *testing.T) {
	setupFixture(t, "fixture/teams.yaml", "fixture/queues.yaml")
	ts := setupPubsub(t)
	defer ts.Close()

	succeedPayload := []byte(`{
        "valid": true,
        "request_count": 3651,
        "elapsed_time": 0,
        "response": {
                "success": 1452,
                "redirect": 2199,
                "client_error": 0,
                "server_error": 0,
                "exception": 0
        },
        "violations": []
}`)
	failedPayload := []byte(`{
        "valid": false,
        "request_count": 100,
        "elapsed_time": 200,
        "response": {
                "success": 98,
                "redirect": 0,
                "client_error": 0,
                "server_error": 2,
                "exception": 0
        },
        "violations": [
                {
                        "request_type": "DUMMY",
                        "description": "アプリケーションが100ミリ秒以内に応答しませんでした",
                        "num": 2
                }
        ]
}`)

	cases := []struct {
		payload []byte
	}{
		{succeedPayload},
		{failedPayload},
	}
	for i, c := range cases {
		targetTeamID := 1
		msgID := publishDummyBenchmarkQueue(t, ts, targetTeamID)
		publishDummyBenchmarkResult(t, ts, targetTeamID, msgID, c.payload)

		q, err := NewQueue(ts.URL)
		if err != nil {
			t.Fatalf("#%d: want non error, got %v", i, err)
		}
		ctx := context.Background()
		err = q.PullAndSave(ctx)
		if err != nil {
			t.Fatalf("#%d: want non error, got %v", i, err)
		}
	}
}

func TestCurrentQueues(t *testing.T) {
	setupFixture(t, "fixture/teams.yaml", "fixture/queues.yaml")
	ts := setupPubsub(t)
	defer ts.Close()

	q, err := NewQueue(ts.URL)
	if err != nil {
		t.Fatalf("want non error, got %v", err)
	}
	publishes := []int{1, 2, 3}
	ctx := context.Background()
	for _, id := range publishes {
		_, err := q.Publish(ctx, id)
		if err != nil {
			t.Fatalf("want non error, got %v", err)
		}
	}

	myTeam := 1
	d, err := NewDatastore()
	if err != nil {
		t.Fatalf("want non error, got %v", err)
	}
	args := make([]interface{}, len(publishes))
	for i, v := range publishes {
		args[i] = v
	}
	rows, err := d.query("select team_id, msg_id from queues where team_id in (?, ?, ?)", args...)
	if err != nil {
		t.Fatalf("want non error, got %v", err)
	}
	expectQueues := []CurrentQueue{}
	for rows.Next() {
		var (
			msgID  string
			teamID int
		)
		err := rows.Scan(&teamID, &msgID)
		if err != nil {
			t.Fatalf("want non error, got %v", err)
		}
		expectQueues = append(expectQueues, CurrentQueue{msgID, teamID == myTeam})
	}

	qs, err := q.CurrentQueues(ctx, myTeam)
	if err != nil {
		t.Fatalf("want non error, got %v", err)
	}
	if !reflect.DeepEqual(expectQueues, qs) {
		t.Errorf("want %v, got %v", expectQueues, qs)
	}
}
