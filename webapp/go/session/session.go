package session

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sync"
	"time"
)

type Manager struct {
	provider    Provider
	cookieName  string
	maxLifeTime int
	lock        sync.Mutex
}

type Provider interface {
	SessionInit(sid string) (Session, error)
	SessionRead(sid string) (Session, error)
	SessionDestroy(sid string) error
}

type Session interface {
	Set(key, value interface{}) error
	Get(key interface{}) interface{}
	Delete(key interface{}) error
	SessionID() string
	AccessedAt() int
}

func NewManager(provideName, cookieName string, maxLifeTime int) (*Manager, error) {
	provider, ok := provides[provideName]
	if !ok {
		return nil, fmt.Errorf("session: unknown provide %q (forgotten import?)", provideName)
	}
	return &Manager{provider: provider, cookieName: cookieName, maxLifeTime: maxLifeTime}, nil
}

var provides = make(map[string]Provider)

func Register(name string, provider Provider) error {
	if provider == nil {
		return fmt.Errorf("Register session provider is nil: %s")
	}

	if _, dup := provides[name]; dup {
		return fmt.Errorf("Already registered %s in provider", name)
	}
	provides[name] = provider
	return nil
}

func (manager *Manager) sessionId() string {
	b := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return ""
	}
	return base64.URLEncoding.EncodeToString(b)
}

func (manager *Manager) SessionStart(w http.ResponseWriter, r *http.Request) (Session, error) {
	manager.lock.Lock()
	defer manager.lock.Unlock()

	cookie, err := r.Cookie(manager.cookieName)
	if err != nil || cookie.Value == "" {
		// Create new cookie
		sid := manager.sessionId()
		session, err := manager.provider.SessionInit(sid)
		if err != nil {
			return nil, fmt.Errorf("SessionInit got invalid response: %v", err)
		}
		cookie := http.Cookie{
			Name:     manager.cookieName,
			Value:    url.QueryEscape(sid),
			Path:     "/",
			HttpOnly: true,
			MaxAge:   manager.maxLifeTime,
		}
		http.SetCookie(w, &cookie)
		return session, nil
	}

	// Use an existing cookie
	sid, _ := url.QueryUnescape(cookie.Value)
	session, err := manager.provider.SessionRead(sid)
	if err != nil {
		return nil, fmt.Errorf("SessionInit got invalid response: %s")
	}
	return session, nil
}

func (manager *Manager) SessionDestroy(w http.ResponseWriter, r *http.Request) error {
	cookie, err := r.Cookie(manager.cookieName)
	if err != nil || cookie.Value == "" {
		return fmt.Errorf("session: invalid cookie")
	}

	manager.lock.Lock()
	defer manager.lock.Unlock()
	manager.provider.SessionDestroy(cookie.Value)
	// update cookie
	cookie.Value = ""
	cookie.Expires = time.Now()
	cookie.MaxAge = -1
	http.SetCookie(w, cookie)
	return nil
}
