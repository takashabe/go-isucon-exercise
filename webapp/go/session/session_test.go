package session

import (
	"fmt"
	"net/http"
	"testing"
)

type TestProvider struct {
}

func NewTestProvider() *TestProvider {
	return &TestProvider{}
}
func (p *TestProvider) SessionInit(sid string) (Session, error) { return nil, nil }
func (p *TestProvider) SessionRead(sid string) (Session, error) { return nil, nil }
func (p *TestProvider) SessionDestroy(sid string) error         { return nil }
func (p *TestProvider) SessionGC(maxLifeTime int64)             {}

type TestResponseWriter struct {
	headers http.Header
	body    []byte
	status  int
}

func NewTestResponseWriter() *TestResponseWriter {
	return &TestResponseWriter{
		headers: make(http.Header),
	}
}
func (w *TestResponseWriter) Header() http.Header {
	return w.headers
}
func (w *TestResponseWriter) Write(body []byte) (int, error) {
	w.body = body
	return len(body), nil
}
func (w *TestResponseWriter) WriteHeader(status int) {
	w.status = status
}

func getProvider(providerName string) (Provider, error) {
	// TODO: switch development and production Provider by env flag
	p := NewTestProvider()
	err := Register(providerName, p)
	if err != nil {
		return nil, fmt.Errorf("getProvider: %s", err)
	}
	return p, nil
}

func getManager(providerName string, cookieName string) (*Manager, error) {
	_, err := getProvider(providerName)
	if err != nil {
		return nil, fmt.Errorf("getProvider: %s")
	}
	m, err := NewManager(providerName, cookieName, 86400)
	if err != nil {
		return nil, fmt.Errorf("getManager: %s")
	}
	return m, nil
}

func TestRegister(t *testing.T) {
	err := Register("test2", nil)
	if err == nil {
		t.Errorf("got nil want error")
	}

	p := NewTestProvider()
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

	providerName := "TestNewManager"
	_, err = getProvider(providerName)
	if err != nil {
		t.Error(err)
	}
	_, err = NewManager(providerName, "go_test", 86400)
	if err != nil {
		t.Errorf("NewManager should be no error, but got error")
	}
}

func TestSessionId(t *testing.T) {
	m, err := getManager("TestSessionId", "gosess")
	if err != nil {
		t.Error(err)
	}
	a := m.sessionId()
	b := m.sessionId()
	if a == b {
		t.Errorf("session IDs should be different, but identical")
	}
}

func TestSessionStart(t *testing.T) {
}
