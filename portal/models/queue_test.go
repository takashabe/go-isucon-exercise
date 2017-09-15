package models

import (
	"net/http/httptest"
	"testing"

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

func TestNewQueue(t *testing.T) {
	ts := setupPubsub(t)
	defer ts.Close()

	_, err := NewQueue(ts.URL)
	if err != nil {
		t.Fatalf("want non error, got %v", err)
	}
}
