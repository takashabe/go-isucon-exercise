package main

import (
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
)

type Helper struct{}

var helper = Helper{}

func (h *Helper) setAddr(ts *httptest.Server, ctx Ctx) Ctx {
	addr := strings.Split(ts.Listener.Addr().String(), ":")
	ctx.host = addr[0]
	ctx.port, _ = strconv.Atoi(addr[1])
	return ctx
}

func TestInitTask(t *testing.T) {
	cases := []struct {
		responseCode int
		expectValid  bool
	}{
		{200, true},
		{500, false},
		{0, false},
	}
	for i, c := range cases {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(c.responseCode)
		}))
		worker := newWorker()
		worker.ctx = helper.setAddr(ts, worker.ctx)
		task := InitTask{
			w: *worker,
		}
		task.Task()

		got := task.FinishHook(*worker.result)
		if got.Valid != c.expectValid {
			t.Errorf("#%d: want: %v, got: %v", i, c.expectValid, got.Valid)
		}
	}
}
