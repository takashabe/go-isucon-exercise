package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

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
		ctx := helper.testCtx(ts)
		driver := helper.testDriver(ctx)

		task := &InitTask{}
		task.Task(ctx, driver)
		got := task.FinishHook(*driver.result)
		if got.Valid != c.expectValid {
			t.Errorf("#%d: want: %v, got: %v", i, c.expectValid, got.Valid)
		}
	}
}
