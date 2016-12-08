package memory

import (
	"container/list"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/takashabe/go-isucon-exercise/webapp/go/session"
)

var pder = &Provider{list: list.New()}

type Provider struct {
	lock     sync.Mutex
	sessions map[string]*list.Element
	list     *list.List
}

type SessionStore struct {
	sid        string
	accessedAt time.Time
	values     map[interface{}]interface{}
}

func (s *SessionStore) Set(key, value interface{}) error {
	s.values[key] = value
	pder.SessionUpdate(s.sid)
	return nil
}

func (s *SessionStore) Get(key interface{}) interface{} {
	pder.SessionUpdate(s.sid)
	if v, ok := s.values[key]; ok {
		return v
	}
	return nil
}

func (s *SessionStore) Delete(key interface{}) error {
	delete(s.values, key)
	pder.SessionUpdate(s.sid)
	return nil
}

func (s *SessionStore) SessionID() string {
	return s.sid
}

func (s *SessionStore) AccessedAt() int {
	return s.accessedAt.Nanosecond()
}

func (s *SessionStore) String() string {
	return fmt.Sprintf("sid: %s, accessedAt: %s, values: %s", s.sid, s.accessedAt, s.values)
}

func (p *Provider) SessionInit(sid string) (session.Session, error) {
	p.lock.Lock()
	defer p.lock.Unlock()
	v := make(map[interface{}]interface{}, 0)
	s := &SessionStore{
		sid:        sid,
		accessedAt: time.Now(),
		values:     v,
	}
	e := p.list.PushBack(s)
	p.sessions[sid] = e
	return s, nil
}

func (p *Provider) SessionRead(sid string) (session.Session, error) {
	if e, ok := p.sessions[sid]; ok {
		return e.Value.(*SessionStore), nil
	}
	s, err := p.SessionInit(sid)
	if err != nil {
		return nil, err
	}
	return s, nil
}

func (p *Provider) SessionDestroy(sid string) error {
	if e, ok := p.sessions[sid]; ok {
		delete(p.sessions, sid)
		p.list.Remove(e)
		return nil
	}
	return nil
}

func (p *Provider) SessionUpdate(sid string) error {
	p.lock.Lock()
	defer p.lock.Unlock()
	if e, ok := p.sessions[sid]; ok {
		e.Value.(*SessionStore).accessedAt = time.Now()
		return nil
	}
	return fmt.Errorf("SessionUpdate: Not found session")
}

func init() {
	fmt.Println("hogehoge")
	pder.sessions = make(map[string]*list.Element)
	err := session.Register("memory", pder)
	if err != nil {
		log.Printf("Failure register memory session provider: %v", err)
	}
}
