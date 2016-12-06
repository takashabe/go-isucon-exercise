package memory

import (
	"testing"
	"time"

	"github.com/takashabe/go-isucon-exercise/webapp/go/session"
)

func getProvider() *Provider {
	return pder
}

func TestSet(t *testing.T) {
	// test set object
	sid := "TestSet"
	session.NewManager("memory", sid, 3600)
	p := getProvider()
	s, _ := p.SessionInit(sid)
	s.Set("test", "value")

	// test update
	beforeTime := s.AccessedAt()
	s.Set("test2", "value")
	afterTime := s.AccessedAt()
	if beforeTime >= afterTime {
		t.Errorf("Not update session, when call Set(): before=%s, after=%s", beforeTime, afterTime)
	}
}

func TestGet(t *testing.T) {
	// test get object
	sid := "TestGet"
	session.NewManager("memory", sid, 3600)
	p := getProvider()
	s, _ := p.SessionInit(sid)
	expected := "value"
	s.Set("test", expected)
	actual := s.Get("test")
	if expected != actual {
		t.Errorf("Different values: expected=%s, actual=%s", expected, actual)
	}

	// test it will be updated when the same key is set
	expected = "value2"
	s.Set("test2", "value")
	s.Set("test2", expected)
	actual = s.Get("test2")
	if expected != actual {
		t.Errorf("Different values has been returned, when updated with set same key: expected=%s, actual=%s", expected, actual)
	}

	// test update session
	beforeTime := s.AccessedAt()
	s.Get("test")
	afterTime := s.AccessedAt()
	if beforeTime >= afterTime {
		t.Errorf("Not update session, when call Get(): before=%s, after=%s", beforeTime, afterTime)
	}

	// test not exist key
	v := s.Get("test999")
	if v != nil {
		t.Errorf("Want nil, but got not nil: actual=%v", v)
	}
}

func TestDelete(t *testing.T) {
	// test delete object
	sid := "TestDelete"
	session.NewManager("memory", sid, 3600)
	p := getProvider()
	s, _ := p.SessionInit(sid)
	s.Set("test", "value")
	s.Delete("test")
	if v := s.Get("test"); v != nil {
		t.Errorf("Want nil, but got value: actual=%v", v)
	}

	// test update session
	beforeTime := s.AccessedAt()
	s.Delete("test")
	afterTime := s.AccessedAt()
	if beforeTime >= afterTime {
		t.Errorf("Not update session, when call Delete(): before=%s, after=%s", beforeTime, afterTime)
	}
}

func TestSessionInit(t *testing.T) {
	sid := "TestSessionInit"
	session.NewManager("memory", sid, 3600)
	p := getProvider()

	beforeTime := time.Now()
	s, err := p.SessionInit(sid)
	if err != nil {
		t.Errorf("Failure initialize session: %s", err)
	}
	afterTime := time.Now()
	accessedAt := s.AccessedAt()
	if beforeTime.Nanosecond() > accessedAt || accessedAt > afterTime.Nanosecond() {
		t.Errorf("Invalid SessionStore accessedAt field")
	}
	if sid != s.SessionID() {
		t.Errorf("Want sid=%s, but got %s", sid, s.SessionID())
	}
}

func TestSessionRead(t *testing.T) {
	// Test if Session exist
	sid := "TestSessionRead1"
	session.NewManager("memory", sid, 3600)
	p := getProvider()
	p.SessionInit(sid)
	s, err := p.SessionRead(sid)
	if err != nil || s == nil {
		t.Errorf("session exist if want session , but got nil")
	}

	// Test if Session and Provider does not exist
	s, err = p.SessionRead("TestSessionRead2")
	if err != nil || s == nil {
		t.Errorf("session does not exist if want session , but got nil")
	}
}

func TestSessionDestroy(t *testing.T) {
	// Test destroy when exist object
	sid := "TestSessionDestroy"
	session.NewManager("memory", sid, 3600)
	p := getProvider()
	p.SessionInit(sid)
	beforeSize := len(p.sessions)
	p.SessionDestroy(sid)
	afterSize := len(p.sessions)
	if beforeSize <= afterSize {
		t.Errorf("Sessions are not decreasing before and after SessionDestroy()")
	}

	// Test destroy when not exist object
	sid2 := "TestSessionDestroy2"
	err := p.SessionDestroy(sid2)
	if err != nil {
		t.Errorf("Want return nil, when not exsit session: actual=%v", err)
	}
}

func TestSessionUpdate(t *testing.T) {
	sid := "TestSessionUpdate"
	session.NewManager("memory", sid, 3600)
	p := getProvider()
	s, _ := p.SessionInit(sid)

	// Test if Session exist
	beforeTime := s.AccessedAt()
	p.SessionUpdate(sid)
	afterTime := s.AccessedAt()
	if beforeTime >= afterTime {
		t.Errorf("session.accessedAt not update if SessionUpdate(): before=%d, after=%d", beforeTime, afterTime)
	}

	// Test if Session does not exist
	err := p.SessionUpdate("TestSessionUpdate_notExist")
	if err == nil {
		t.Errorf("want nil, but got session object")
	}
}
