package main

import (
	"reflect"
	"testing"

	"github.com/pkg/errors"
)

func TestCreateSessions(t *testing.T) {
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
		m := Master{}
		got, err := m.createSessions(c.input)
		if errors.Cause(err) != c.expectErr {
			t.Errorf("#%d: want: %v, got: %v", i, c.expectErr, err)
		}
		if !reflect.DeepEqual(got, c.expectObject) {
			t.Errorf("#%d: want: %v, got: %v", i, c.expectObject, got)
		}
	}
}
