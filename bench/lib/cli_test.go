package bench

import (
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/pkg/errors"
)

func TestParseArgs(t *testing.T) {
	cli := &CLI{
		outStream: os.Stdout,
		errStream: os.Stderr,
	}
	cases := []struct {
		input       string
		expectErr   error
		expectParam param
	}{
		{
			"",
			nil,
			param{
				file: defaultFile,
				host: defaultHost,
				port: defaultPort,
				time: defaultTime,
			},
		},
		{
			"-port 8080 -host 127.0.0.1 -file data/test.json",
			nil,
			param{
				host: "127.0.0.1",
				port: 8080,
				time: defaultTime,
				file: "data/test.json",
			},
		},
		{
			"-nonparam foo",
			ErrParseFailed,
			param{},
		},
	}
	for i, c := range cases {
		p := &param{}
		err := cli.parseArgs(strings.Split(c.input, " "), p)
		if errors.Cause(err) != c.expectErr {
			t.Errorf("#%d: want %#v, got %#v", i, c.expectErr, err)
		}
		if reflect.DeepEqual(p, c.expectParam) {
			t.Errorf("#%d: want %d, got %d", i, c.expectParam, p)
		}
	}
}
