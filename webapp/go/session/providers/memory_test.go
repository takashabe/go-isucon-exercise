package memory

import (
	"fmt"
	"testing"

	"github.com/takashabe/go-isucon-exercise/webapp/go/session"
)

func TestSessionInit(t *testing.T) {
	m, err := session.NewManager("memory", "gosessid", 3600)
	if err != nil {
		t.Errorf("Failure create session manager: %s", err)
	}
	fmt.Println(m.provider)
}
