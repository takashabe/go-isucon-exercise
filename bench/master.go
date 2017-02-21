package main

import (
	"encoding/json"
	"io/ioutil"

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

type Master struct {
}

func (m *Master) start() {
}

func (m *Master) createSessions(path string) (*UserSchemas, error) {
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
