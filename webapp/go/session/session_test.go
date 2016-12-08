package session

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
	"time"
)

// DummyProvider is testing provider. provider is expected implement to provider package
type DummyProvider struct{}

func (p *DummyProvider) SessionInit(sid string) (Session, error) { return NewDummySession(), nil }
func (p *DummyProvider) SessionRead(sid string) (Session, error) { return NewDummySession(), nil }
func (p *DummyProvider) SessionDestroy(sid string) error         { return nil }
func (p *DummyProvider) SessionGC(maxLifeTime int64)             {}

func NewDummyProvider() *DummyProvider {
	return &DummyProvider{}
}

// DummySession is testing session. session is expected implemnt to provider package
type DummySession struct{}

func (s *DummySession) Set(key, value interface{}) error { return nil }
func (s *DummySession) Get(key interface{}) interface{}  { return nil }
func (s *DummySession) Delete(key interface{}) error     { return nil }
func (s *DummySession) SessionID() string                { return "" }
func (s *DummySession) AccessedAt() int                  { return 0 }

func NewDummySession() *DummySession {
	return &DummySession{}
}

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

func getProvider(providerName string) Provider {
	// TODO: switch development and production Provider by env flag
	p := NewDummyProvider()
	Register(providerName, p)
	return p
}

func getManager(providerName string, cookieName string) (*Manager, error) {
	getProvider(providerName)
	m, err := NewManager(providerName, cookieName, 3600)
	if err != nil {
		return nil, fmt.Errorf("getManager: %s")
	}
	return m, nil
}

func TestRegister(t *testing.T) {
	// Test exist provider
	p := NewDummyProvider()
	err := Register("test", p)
	if err != nil {
		t.Errorf("got error want nil")
	}

	// Test not exsit provider
	err = Register("test2", nil)
	if err == nil {
		t.Errorf("Want err, but got not error")
	}

	// Test duplicate provider
	Register("test3", p)
	err = Register("test3", p)
	if err == nil {
		t.Errorf("Want register error, when already same provider name")
	}
}

func TestNewManager(t *testing.T) {
	_, err := NewManager("", "go_test", 86400)
	if err == nil {
		t.Errorf("NewManager invalid argument so expected got error, but no error")
	}

	providerName := "TestNewManager"
	getProvider(providerName)
	_, err = NewManager(providerName, "go_test", 86400)
	if err != nil {
		t.Errorf("NewManager should be no error, but got error")
	}
}

func TestSessionId(t *testing.T) {
	m, _ := getManager("TestSessionId", "gosess")
	a := m.sessionId()
	b := m.sessionId()
	if a == b {
		t.Errorf("session IDs should be different, but identical")
	}
}

func TestSessionStart(t *testing.T) {
	// Test set session(cookie)
	res := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/", nil)
	m, _ := getManager("TestSessionStart", "gosess")
	if s, err := m.SessionStart(res, req); s == nil || err != nil {
		t.Errorf("Want have return session and not error: actual session=%v, error=%v", s, err)
	}
	_, err := getCookieValue(res.Header(), m.cookieName)
	if err != nil {
		t.Errorf("Invalid cookie, not found cookie name: %s", m.cookieName)
	}

	actualAge, _ := getCookieValue(res.Header(), "Max-Age")
	if a, _ := strconv.Atoi(actualAge); a != m.maxLifeTime {
		t.Errorf("Invalid cookie, Max-Age want %s but got %s", m.maxLifeTime, actualAge)
	}

	// Test already existing cookie
	// set dummy cookie
	sid, _ := getCookieValue(res.Header(), m.cookieName)
	cookie := http.Cookie{
		Name:   m.cookieName,
		Value:  sid,
		MaxAge: m.maxLifeTime,
	}
	req.AddCookie(&cookie)
	if s, err := m.SessionStart(res, req); s == nil || err != nil {
		t.Errorf("Want have return session and not error: actual session=%v, error=%v", s, err)
	}
}

func TestSessionDestroy(t *testing.T) {
	res := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/", nil)
	m, _ := getManager("TestSessionDestroy", "gosess")
	m.SessionStart(res, req)

	// set dummy cookie
	sid, _ := getCookieValue(res.Header(), m.cookieName)
	cookie := http.Cookie{
		Name:   m.cookieName,
		Value:  sid,
		MaxAge: m.maxLifeTime,
	}
	res.Header().Del("Set-Cookie")
	req.AddCookie(&cookie)

	if err := m.SessionDestroy(res, req); err != nil {
		t.Errorf("Want non error, but got error: err=%v", err)
	}
	maxAge, _ := getCookieValue(res.Header(), "Max-Age")
	if i, _ := strconv.Atoi(maxAge); i > 0 {
		// "MaxAge<0" replace by "MaxAge=0"
		t.Errorf("MaxAge want 0, but got MaxAge=%s", maxAge)
	}
	expires, _ := getCookieValue(res.Header(), "Expires")
	// Expires are saved in seconds. So only to test what is likely.
	actualTime, _ := time.Parse(time.RFC1123, expires)
	expectedTime := time.Now().Add(time.Duration(10) * time.Second)
	if expectedTime.Before(actualTime) {
		t.Errorf("Expires are not saved time.Now()")
	}
}

func getCookieValue(h http.Header, filter string) (string, error) {
	cookie := h.Get("Set-Cookie")
	parts := strings.Split(strings.TrimSpace(cookie), ";")
	for _, v := range parts {
		v = strings.TrimSpace(v)
		exist := strings.HasPrefix(v, filter)
		if !exist {
			continue
		}
		// trim string (filter + "=")
		return v[len(filter)+1:], nil
	}
	return "", fmt.Errorf("Not found value")
}
