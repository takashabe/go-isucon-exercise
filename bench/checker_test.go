package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func testChecker() *Checker {
	return &Checker{
		ctx:         defaultCtx,
		result:      newResult(),
		path:        "/",
		requestName: "TEST",
		response:    *httptest.NewRecorder().Result(),
	}
}

func testResponse(code int) *http.Response {
	recorder := httptest.NewRecorder()
	recorder.WriteHeader(code)
	return recorder.Result()
}

func testRedirectResponse(path string, code int) *http.Response {
	recorder := httptest.NewRecorder()
	recorder.Header().Set("Location", path)
	recorder.WriteHeader(code)
	return recorder.Result()
}

func TestIsStatusCode(t *testing.T) {
	cases := []struct {
		input        int
		expectResult *Result
	}{
		{
			200,
			newResult(),
		},
		{
			500,
			newResult().addViolation("TEST", fmt.Sprintf(causeStatusCode, "/", 500, 200)),
		},
	}
	for i, c := range cases {
		checker := testChecker()
		checker.isStatusCode(c.input)
		if !reflect.DeepEqual(checker.result, c.expectResult) {
			t.Errorf("#%d: want %v, got %v", i, c.expectResult, checker.result)
		}
	}
}

func TestIsRedirect(t *testing.T) {
	cases := []struct {
		response     *http.Response
		input        string
		expectResult *Result
	}{
		{
			testRedirectResponse("/test", 302),
			"/test",
			newResult(),
		},
		{
			testRedirectResponse("/test", 200),
			"/test",
			newResult().addViolation("TEST", fmt.Sprintf(causeRedirectStatusCode, 200)),
		},
		{
			testResponse(302),
			"/test",
			newResult().addViolation("TEST", fmt.Sprintf(causeNoneLocation)),
		},
		{
			testRedirectResponse("http://localhost/test", 302),
			"/test",
			newResult(),
		},
		{
			testRedirectResponse("http://localhost/foo", 302),
			"/test",
			newResult().addViolation("TEST", fmt.Sprintf(causeInvalidLocationPath, "/test", "/foo")),
		},
	}
	for i, c := range cases {
		checker := testChecker()
		checker.response = *c.response
		checker.isRedirect(c.input)
		if !reflect.DeepEqual(checker.result, c.expectResult) {
			t.Errorf("#%d: want %v, got %v", i, c.expectResult, checker.result)
		}
	}
}
