package agent

import (
	"context"
	"net/http/httptest"
	"net/url"
	"reflect"
	"testing"
	"time"

	"github.com/takashabe/go-pubsub/client"
	"github.com/takashabe/go-pubsub/server"
)

func setupPubsubServer(t *testing.T) *httptest.Server {
	s, err := server.NewServer("testdata/config.yml")
	if err != nil {
		t.Fatalf("want non error, got %v", err)
	}
	if err := s.PrepareServer(); err != nil {
		t.Fatalf("want non error, got %v", err)
	}
	ts := httptest.NewServer(server.Routes())
	for _, n := range []string{pullServerName, publishServerName} {
		topic := createTopic(t, ts, n)
		createSubscription(t, ts, n, topic)
	}
	return ts
}

func createTopic(t *testing.T, ts *httptest.Server, id string) *client.Topic {
	ctx := context.Background()
	client, err := client.NewClient(ctx, ts.URL)
	if err != nil {
		t.Fatalf("want non error, got %v", err)
	}
	topic := client.Topic(id)
	exist, err := topic.Exists(ctx)
	if err != nil {
		t.Fatalf("want non error, got %v", err)
	}
	if exist {
		return topic
	}

	_, err = client.CreateTopic(ctx, id)
	if err != nil {
		t.Fatalf("want non error, got %v", err)
	}
	return topic
}

func createSubscription(t *testing.T, ts *httptest.Server, id string, topic *client.Topic) {
	ctx := context.Background()
	c, err := client.NewClient(ctx, ts.URL)
	if err != nil {
		t.Fatalf("want non error, got %v", err)
	}
	exist, err := c.Subscription(id).Exists(ctx)
	if err != nil {
		t.Fatalf("want non error, got %v", err)
	}
	if exist {
		return
	}
	cfg := client.SubscriptionConfig{Topic: topic}
	if _, err := c.CreateSubscription(ctx, id, cfg); err != nil {
		t.Fatalf("want non error, got %v", err)
	}
}

func publishDummyBenchmarkRequest(t *testing.T, ts *httptest.Server) string {
	ctx := context.Background()
	c, err := client.NewClient(ctx, ts.URL)
	if err != nil {
		t.Fatalf("want non error, got %v", err)
	}
	result := c.Topic(pullServerName).Publish(ctx, &client.Message{
		Attributes: map[string]string{"team_id": "1"},
	})
	msgID, err := result.Get(ctx)
	if err != nil {
		t.Fatalf("want non error, got %v", err)
	}
	return msgID
}

func TestNewAgent(t *testing.T) {
	ts := setupPubsubServer(t)
	defer ts.Close()

	cases := []struct {
		inputInterval int
		inputPubsub   string
		expectErr     error
	}{
		{0, ts.URL, nil},
		{0, "invalidURL", &url.Error{}},
	}
	for i, c := range cases {
		d, err := NewDispatch("./testdata/dummyScript", "./testdata/dummyParam", "localhost", 80)
		if err != nil {
			t.Fatalf("#%d: want non error, got %v", i, err)
		}
		_, err = NewAgent(c.inputInterval, c.inputPubsub, d)
		if reflect.TypeOf(err) != reflect.TypeOf(c.expectErr) {
			t.Errorf("#%d: want error %v, got %v", i, c.expectErr, err)
		}
	}
}

func TestPolling(t *testing.T) {
	ts := setupPubsubServer(t)
	defer ts.Close()

	id := publishDummyBenchmarkRequest(t, ts)
	ctx, cancel := context.WithCancel(context.Background())

	d, err := NewDispatch("./testdata/dummyScript", "./testdata/dummyParam", "localhost", 80)
	if err != nil {
		t.Fatalf("want non error, got %v", err)
	}
	agent, err := NewAgent(0, ts.URL, d)
	if err != nil {
		t.Fatalf("want non error, got %v", err)
	}
	agent.interval = 20 * time.Millisecond
	msg, err := agent.Polling(ctx)
	if err != nil {
		t.Fatalf("want non error, got %v", err)
	}
	if msg.ID != id {
		t.Errorf("want message ID %s, got %s", id, msg.ID)
	}

	// expect not found message
	go func() {
		_, err = agent.Polling(ctx)
		if err != context.Canceled {
			t.Errorf("want error %v, got %v", context.Canceled, err)
		}
	}()
	time.Sleep(40 * time.Millisecond)
	cancel()
}

func TestDispatch(t *testing.T) {
	ts := setupPubsubServer(t)
	defer ts.Close()

	d, err := NewDispatch("./testdata/dummyScript", "./testdata/dummyParam", "localhost", 80)
	if err != nil {
		t.Fatalf("want non error, got %v", err)
	}
	agent, err := NewAgent(0, ts.URL, d)
	if err != nil {
		t.Fatalf("want non error, got %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	_, err = agent.Dispatch(ctx)
	if err != nil {
		t.Fatalf("want non error, got %v", err)
	}

	agent.dispatch.script = "./testdata/sleepScript"
	go func() {
		_, err := agent.Dispatch(ctx)
		if err != context.Canceled {
			t.Errorf("want error %v, got %v", context.Canceled, err)
		}
	}()
	cancel()
}

func TestSendResult(t *testing.T) {
	ts := setupPubsubServer(t)
	defer ts.Close()

	d, err := NewDispatch("./testdata/dummyScript", "./testdata/dummyParam", "localhost", 80)
	if err != nil {
		t.Fatalf("want non error, got %v", err)
	}
	agent, err := NewAgent(0, ts.URL, d)
	if err != nil {
		t.Fatalf("want non error, got %v", err)
	}

	ctx := context.Background()
	err = agent.SendResult(ctx, []byte("dummy"), map[string]string{"dummy": "dummy"})
	if err != nil {
		t.Fatalf("want non error, got %v", err)
	}
}
