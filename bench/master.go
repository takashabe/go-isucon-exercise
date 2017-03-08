package main

import (
	"encoding/json"
	"io/ioutil"
	"reflect"

	"github.com/pkg/errors"
)

var (
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

type Master struct{}

func (m *Master) start(file, host string, port, time int) ([]byte, error) {
	// TODO
	// 1. create workers
	// 2. run for each workers with order()
	// 3. sum return results from worker.run

	// TODO: export run class parameter. for example param.json
	session, err := m.getSessions(file)
	if err != nil {
		return nil, err
	}
	result := newResult()
	workers := IsuconWorkers()
	for _, w := range workers {
		w.ctx.host = host
		w.ctx.port = port
		for _, t := range w.tasks {
			PrintDebugf("RUN %s", reflect.ValueOf(t).String())
			t.SetWorker(*w)
			t.Task(session)
			r := t.FinishHook()
			result = result.Merge(r)
			if !result.Valid {
				PrintDebugf("invalid result: %#v\n", t)
				break
			}
		}
	}

	json, err := result.json()
	if err != nil {
		PrintDebugf("failed to result.json(): %s", err.Error())
		return nil, err
	}
	return json, nil
}

var sessions []*Session = nil

func (m *Master) getSessions(path string) ([]*Session, error) {
	if sessions != nil {
		return sessions, nil
	}

	p, err := m.loadParams(path)
	if err != nil {
		return nil, err
	}
	sessions = make([]*Session, len(p.Parameters))
	for i, v := range p.Parameters {
		sessions[i] = newSession(v)
	}
	return sessions, nil
}

func (m *Master) loadParams(path string) (*UserSchemas, error) {
	data, err := ioutil.ReadFile(path)
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
