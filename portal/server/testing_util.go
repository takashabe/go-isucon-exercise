package server

import (
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"strings"
	"testing"

	fixture "github.com/takashabe/go-fixture"
	_ "github.com/takashabe/go-fixture/mysql" // mysql driver
	"github.com/takashabe/go-isucon-exercise/portal/models"
	"github.com/takashabe/go-message-queue/server"
)

func clientWithNonRedirect() *http.Client {
	// suppression to redirect
	return &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
}

func clientWithCookie(t *testing.T) *http.Client {
	jar, err := cookiejar.New(nil)
	if err != nil {
		t.Fatalf("want non error, got %v", err)
	}
	return &http.Client{
		Jar: jar,
	}
}

func setupFixture(t *testing.T, files ...string) {
	db, err := models.NewDatastore()
	if err != nil {
		t.Fatalf("want non nil, got %v", err)
	}
	f := fixture.NewFixture(db.Conn, "mysql")
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

func setupServer(t *testing.T, pubsubAddr string) *httptest.Server {
	server, err := NewServer(pubsubAddr)
	if err != nil {
		t.Fatalf("want non error, got %v", err)
	}
	return httptest.NewServer(server.Routes())
}

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
