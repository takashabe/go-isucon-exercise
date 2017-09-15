package models

import (
	"context"
	"fmt"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/takashabe/go-message-queue/client"
	"github.com/takashabe/go-message-queue/server"
)

func setupPubsub(t *testing.T) *httptest.Server {
	s, err := server.NewServer("testdata/config.yaml")
	if err != nil {
		t.Fatalf("failed to server.NewServer, error=%v", err)
	}
	if err := s.InitDatastore(); err != nil {
		t.Fatalf("failed to server.InitDatastore, error=%v", err)
	}
	return httptest.NewServer(server.Routes())
}

func publishDummyBenchmarkResult(t *testing.T, ts *httptest.Server, payload []byte) {
	q, err := NewQueue(ts.URL)
	if err != nil {
		t.Fatalf("want non error, got %v", err)
	}
	ctx := context.Background()
	res := q.c.Topic(PubsubServerName).Publish(ctx, &client.Message{
		Data: payload,
		Attributes: map[string]string{
			"team_id":    "1",
			"created_at": fmt.Sprintf("%d", time.Now().Unix()),
		},
	})
	_, err = res.Get(ctx)
	if err != nil {
		t.Fatalf("want non error, got %v", err)
	}
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
	ts := setupPubsub(t)
	defer ts.Close()

	q, err := NewQueue(ts.URL)
	if err != nil {
		t.Fatalf("want non error, got %v", err)
	}
	ctx := context.Background()
	err = q.Publish(ctx, 1)
	if err != nil {
		t.Fatalf("want non error, got %v", err)
	}
}

func TestPull(t *testing.T) {
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
		publishDummyBenchmarkResult(t, ts, c.payload)
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
