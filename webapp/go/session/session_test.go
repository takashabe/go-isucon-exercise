package session

import "testing"

func TestRegister(t *testing.T) {
	err := Register("test", nil)
	if err == nil {
		t.Errorf("got nil want error")
	}
}

func TestNewManager(t *testing.T) {
}
