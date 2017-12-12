package models

import (
	"net/http/httptest"
	"strings"
	"testing"

	fixture "github.com/takashabe/go-fixture"
	_ "github.com/takashabe/go-fixture/mysql" // mysql driver
	"github.com/takashabe/go-pubsub/server"
)

func setupFixture(t *testing.T, files ...string) {
	db, err := NewDatastore()
	if err != nil {
		t.Fatalf("want non nil, got %v", err)
	}
	f, err := fixture.NewFixture(db.Conn, "mysql")
	if err != nil {
		t.Fatalf("want non nil, got %v", err)
	}
	// always necessary file
	err = f.LoadSQL("fixture/schema.sql")
	if err != nil {
		t.Fatalf("want non error, got %v", err)
	}

	for _, file := range files {
		if strings.HasSuffix(file, ".sql") {
			err = f.LoadSQL(file)
		} else {
			err = f.Load(file)
		}
	}
	if err != nil {
		t.Fatalf("want non nil, got %v", err)
	}
}

func setupPubsub(t *testing.T) *httptest.Server {
	s, err := server.NewServer("testdata/config.yaml")
	if err != nil {
		t.Fatalf("failed to server.NewServer, error=%v", err)
	}
	if err := s.PrepareServer(); err != nil {
		t.Fatalf("failed to server.PrepareServer, error=%v", err)
	}
	return httptest.NewServer(server.Routes())
}
