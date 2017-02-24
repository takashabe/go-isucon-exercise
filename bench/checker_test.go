package main

import (
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestIsStatusCode(t *testing.T) {
	checker := Checker{
		ctx:         defaultCtx,
		result:      newResult(),
		path:        "/",
		requestName: "TEST",
		response:    *httptest.NewRecorder().Result(),
	}
	cases := []struct {
		checker      Checker
		input        int
		expectResult *Result
	}{
		{
			checker,
			200,
			newResult(),
		},
		{
			checker,
			500,
			newResult().addViolation("TEST", "パス '/' へのレスポンスコード 500 が期待されていましたが 200 でした"),
		},
	}
	for i, c := range cases {
		c.checker.isStatusCode(c.input)
		if !reflect.DeepEqual(c.checker.result, c.expectResult) {
			t.Errorf("#%d: want %v, got %v", i, c.expectResult, c.checker.result)
		}
	}
}
