package benchmark

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestBootstrapTask(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(500)
	}))

	ctx := helper.testCtx(ts)
	driver := helper.testDriver(ctx)
	task := &BootstrapTask{}
	task.Task(ctx, driver)

	got := task.FinishHook(*driver.result)
	want := false
	if got.Valid != want {
		t.Errorf("want: %d, got: %d", want, got.Valid)
	}
}
