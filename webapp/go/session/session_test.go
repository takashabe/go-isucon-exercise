package session

import "testing"

type TestProvider struct {
}

func (p *TestProvider) SessionInit(sid string) (Session, error) { return nil, nil }
func (p *TestProvider) SessionRead(sid string) (Session, error) { return nil, nil }
func (p *TestProvider) SessionDestroy(sid string) error         { return nil }
func (p *TestProvider) SessionGC(maxLifeTime int64)             {}

func TestRegister(t *testing.T) {
	err := Register("test2", nil)
	if err == nil {
		t.Errorf("got nil want error")
	}

	p := &TestProvider{}
	err = Register("test", p)
	if err != nil {
		t.Errorf("got error want nil")
	}
}

func TestNewManager(t *testing.T) {
	_, err := NewManager("", "go_test", 86400)
	if err == nil {
		t.Errorf("NewManager invalid argument so expected got error, but no error")
	}

	p := &TestProvider{}
	Register("memory", p)
	m, err := NewManager("memory", "go_test", 86400)
	if err != nil {
		t.Errorf("NewManager should be no error, but got error")
	}
	t.Log(m)
}

func TestSessionId(t *testing.T) {
	// a := SessionId()
	// b := SessionId()
	// if a == b {
	//   t.Errorf("session IDs should be different, but identical")
	// }
}
