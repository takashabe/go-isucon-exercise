package benchmark

import (
	"bytes"
	"reflect"
	"strings"
	"testing"

	"github.com/pkg/errors"
)

func TestParseArgs(t *testing.T) {
	cli := &CLI{
		outStream: new(bytes.Buffer),
		errStream: new(bytes.Buffer),
	}
	cases := []struct {
		input       string
		expectErr   error
		expectParam *param
	}{
		{
			"./bench",
			nil,
			&param{
				host:  defaultHost,
				port:  defaultPort,
				file:  defaultFile,
				agent: defaultAgent,
			},
		},
		{
			"./bench -port 8080 -host 127.0.0.1 -file data/test.json -agent test",
			nil,
			&param{
				host:  "127.0.0.1",
				port:  8080,
				file:  "data/test.json",
				agent: "test",
			},
		},
		{
			"./bench -nonparam foo",
			ErrParseFailed,
			&param{
				host:  defaultHost,
				port:  defaultPort,
				file:  defaultFile,
				agent: defaultAgent,
			},
		},
	}
	for i, c := range cases {
		p := &param{}
		args := strings.Split(c.input, " ")
		err := cli.parseArgs(args[1:], p)
		if errors.Cause(err) != c.expectErr {
			t.Errorf("#%d: want %#v, got %#v", i, c.expectErr, err)
		}
		if !reflect.DeepEqual(p, c.expectParam) {
			t.Errorf("#%d: want %d, got %d", i, c.expectParam, p)
		}
	}
}
