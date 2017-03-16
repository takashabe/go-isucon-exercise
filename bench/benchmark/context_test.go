package main

import (
	"reflect"
	"testing"

	"github.com/pkg/errors"
)

func testSession(p UserSchema) *Session {
	s, _ := newSession(p)
	return s
}

func TestLoadParam(t *testing.T) {
	cases := []struct {
		input        string
		expectObject *UserSchemas
		expectErr    error
	}{
		{
			"testdata/foo.json",
			&UserSchemas{[]UserSchema{
				UserSchema{ID: 1, Name: "a", Email: "b", Password: "c"},
				UserSchema{ID: 2, Name: "d", Email: "e", Password: "f"},
			}},
			nil,
		},
		{
			"testdata/none.json",
			nil,
			ErrFailedReadFile,
		},
		{
			"testdata/invalid.json",
			nil,
			ErrFailedParseJson,
		},
	}
	for i, c := range cases {
		ctx := newCtx()
		ctx.paramFile = c.input
		got, err := ctx.loadParams()
		if errors.Cause(err) != c.expectErr {
			t.Errorf("#%d: want %#v, got %#v", i, c.expectErr, err)
		}
		if !reflect.DeepEqual(c.expectObject, got) {
			t.Errorf("#%d: want %#v, got %#v", i, c.expectObject, got)
		}
	}
}

func TestSetupSessions(t *testing.T) {
	cases := []struct {
		input        string
		expectObject []*Session
		expectErr    error
	}{
		{
			"testdata/foo.json",
			[]*Session{
				testSession(UserSchema{ID: 1, Name: "a", Email: "b", Password: "c"}),
				testSession(UserSchema{ID: 2, Name: "d", Email: "e", Password: "f"}),
			},
			nil,
		},
	}
	for i, c := range cases {
		ctx := newCtx()
		ctx.paramFile = c.input
		err := ctx.setupSessions()
		if err != c.expectErr {
			t.Errorf("#%d: want %#v, got %#v", i, c.expectErr, err)
		}
		act := ctx.sessions
		if !reflect.DeepEqual(c.expectObject, act) {
			t.Errorf("#%d: want %#v, got %#v", i, c.expectObject, act)
		}
	}
}
