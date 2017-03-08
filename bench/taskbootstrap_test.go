package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
)

func testSessions() []*Session {
	m := Master{}
	params, err := m.loadParams("testdata/param.json")
	if err != nil {
		panic(fmt.Sprintf("failed create sessions: %s", err.Error()))
	}

	log.Println(len(params.Parameters))
	sessions := make([]*Session, len(params.Parameters))
	for i, v := range params.Parameters {
		sessions[i] = newSession(v)
	}
	return sessions
}

func setAddr(ts *httptest.Server, ctx Ctx) Ctx {
	addr := strings.Split(ts.Listener.Addr().String(), ":")
	ctx.host = addr[0]
	ctx.port, _ = strconv.Atoi(addr[1])
	return ctx
}

func TestTask(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(500)
	}))
	w := NewWorker()
	w.ctx = setAddr(ts, w.ctx)
	task := BootstrapTask{w: *w}
	task.Task(testSessions())

	want := false
	got := task.FinishHook()
	if got.Valid != want {
		t.Errorf("want: %d, got: %d", want, got.Valid)
	}
}
