package benchmark

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"sync"
	"time"

	"github.com/pkg/errors"
)

var (
	// failed create session errors
	ErrNotFoundSession = errors.New("not found session")
	ErrFailedReadFile  = errors.New("failed to read file")
	ErrFailedParseJson = errors.New("failed to parse json")
)

// UserSchema represents the user column userd in the request
type UserSchemas struct {
	Parameters []UserSchema
}
type UserSchema struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Ctx is environment settings in each Worker
type Ctx struct {
	// url query parameters
	scheme string
	host   string
	port   int
	agent  string

	// request timeout
	getTimeout  time.Duration
	postTimeout time.Duration

	// worker running time
	workerRunningTime time.Duration

	// parameter json file
	paramFile string

	// session list
	sessions []*Session
	mu       sync.Mutex
}

var defaultCtx = Ctx{
	scheme:            "http",
	host:              defaultHost,
	port:              defaultPort,
	agent:             defaultAgent,
	getTimeout:        30 * time.Second,
	postTimeout:       30 * time.Second,
	workerRunningTime: 30 * time.Second,
	paramFile:         defaultFile,
}

func newCtx() *Ctx {
	ctx := defaultCtx
	return &ctx
}

func (c *Ctx) setSessions(sessions []*Session) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.sessions = sessions
}

func (c *Ctx) getSession(i int) (*Session, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if len(c.sessions) > i {
		return nil, ErrNotFoundSession
	}
	return c.sessions[i], nil
}

// create sessions and setting on Ctx
func (c *Ctx) setupSessions() error {
	if c.sessions != nil {
		return nil
	}

	p, err := c.loadParams()
	if err != nil {
		return err
	}
	sessions := make([]*Session, len(p.Parameters))
	for i, v := range p.Parameters {
		s, err := newSession(v)
		if err != nil {
			return err
		}
		sessions[i] = s
	}
	c.sessions = sessions
	return nil
}

func (c *Ctx) loadParams() (*UserSchemas, error) {
	data, err := ioutil.ReadFile(c.paramFile)
	if err != nil {
		return nil, errors.Wrap(ErrFailedReadFile, err.Error())
	}
	var schemas UserSchemas
	err = json.Unmarshal(data, &schemas)
	if err != nil {
		return nil, errors.Wrap(ErrFailedParseJson, err.Error())
	}
	return &schemas, nil
}

func (c *Ctx) uri(path string) string {
	return fmt.Sprintf("%s://%s:%d%s", c.scheme, c.host, c.port, path)
}
