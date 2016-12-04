package memory

import (
	"container/list"
	"testing"
	"time"

	"github.com/takashabe/go-isucon-exercise/webapp/go/session"
)

func getProvider() *Provider {
	return &Provider{list: list.New()}
}

func TestSessionInit(t *testing.T) {
	_, err := session.NewManager("memory", "gosessid", 3600)
	if err != nil {
		t.Errorf("Failure create session manager: %s", err)
	}

	beforeTime := time.Now()
	p := getProvider()
	s, err := p.SessionInit("gosessid")
	if err != nil {
		t.Errorf("Failure initialize session: %s", err)
	}
	afterTime := time.Now()
	accessedAt := s.AccessedAt()
	if beforeTime.Nanosecond() > accessedAt || accessedAt > afterTime.Nanosecond() {
		t.Errorf("Invalid SessionStore accessedAt field")
	}
}

func TestSessionRead(t *testing.T) {
}

func TestSessionDestroy(t *testing.T) {
}

func TestSessionUpdate(t *testing.T) {
}
