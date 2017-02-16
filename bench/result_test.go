package main

import (
	"reflect"
	"testing"
)

func TestAddResponse(t *testing.T) {
	cases := []struct {
		input  int
		expect *Result
	}{
		{
			202,
			&Result{
				requestCount: 1,
				response:     &ResponseCounter{success: 1},
				violations:   make([]*Violation, 0),
			},
		},
		{
			302,
			&Result{
				requestCount: 1,
				response:     &ResponseCounter{redirect: 1},
				violations:   make([]*Violation, 0),
			},
		},
		{
			404,
			&Result{
				requestCount: 1,
				response:     &ResponseCounter{clientError: 1},
				violations:   make([]*Violation, 0),
			},
		},
		{
			999,
			&Result{
				requestCount: 1,
				response:     &ResponseCounter{serverError: 1},
				violations:   make([]*Violation, 0),
			},
		},
		{
			0,
			&Result{
				requestCount: 1,
				response:     &ResponseCounter{serverError: 1},
				violations:   make([]*Violation, 0),
			},
		},
	}
	for i, c := range cases {
		got := newResult().addResponse(c.input)
		if !reflect.DeepEqual(got, c.expect) {
			t.Errorf("#%d: want: %v, got: %v", i, c.expect, got)
		}
	}
}

func TestAddResponseException(t *testing.T) {
	cases := []struct {
		base   *Result
		expect *Result
	}{
		{
			&Result{
				requestCount: 1,
				response:     &ResponseCounter{success: 200},
			},
			&Result{
				requestCount: 2,
				response:     &ResponseCounter{success: 200, exception: 1},
			},
		},
	}
	for i, c := range cases {
		got := c.base.addResponseException()
		if !reflect.DeepEqual(got, c.expect) {
			t.Errorf("#%d: want: %v, got: %v", i, c.expect, got)
		}
	}
}

func TestAddViolation(t *testing.T) {
	cases := []struct {
		inputName  string
		inputCause string
		base       *Result
		expect     *Result
	}{
		{
			"foo",
			"bar",
			&Result{},
			&Result{
				violations: []*Violation{&Violation{
					requestName: "foo",
					cause:       "bar",
					count:       1,
				}},
			},
		},
		{
			"foo",
			"bar",
			&Result{
				violations: []*Violation{
					&Violation{
						requestName: "foo",
						cause:       "bar",
						count:       1,
					},
					&Violation{
						requestName: "hoge",
						cause:       "piyo",
						count:       2,
					},
				},
			},
			&Result{
				violations: []*Violation{
					&Violation{
						requestName: "foo",
						cause:       "bar",
						count:       2,
					},
					&Violation{
						requestName: "hoge",
						cause:       "piyo",
						count:       2,
					},
				},
			},
		},
	}
	for i, c := range cases {
		got := c.base.addViolation(c.inputName, c.inputCause)
		if !reflect.DeepEqual(got, c.expect) {
			t.Errorf("#%d: want: %v, got: %v", i, c.expect, got)
		}
	}
}

func TestMerge(t *testing.T) {
	cases := []struct {
		base   *Result
		input  Result
		expect *Result
	}{
		{
			&Result{
				valid:        false,
				requestCount: 1,
				elapsedTime:  300,
				response:     &ResponseCounter{success: 1, exception: 1},
				violations: []*Violation{
					&Violation{
						requestName: "foo",
						cause:       "bar",
						count:       1,
					},
				},
			},
			Result{
				valid:        true,
				requestCount: 1,
				elapsedTime:  300,
				response:     &ResponseCounter{success: 1},
				violations: []*Violation{
					&Violation{
						requestName: "foo",
						cause:       "bar",
						count:       1,
					},
					&Violation{
						requestName: "hoge",
						cause:       "piyo",
						count:       1,
					},
				},
			},
			&Result{
				valid:        false,
				requestCount: 2,
				elapsedTime:  600,
				response:     &ResponseCounter{success: 2, exception: 1},
				violations: []*Violation{
					&Violation{
						requestName: "foo",
						cause:       "bar",
						count:       2,
					},
					&Violation{
						requestName: "hoge",
						cause:       "piyo",
						count:       1,
					},
				},
			},
		},
		{
			&Result{
				response: &ResponseCounter{success: 1},
			},
			Result{
				response: newResponse(),
			},
			&Result{
				response: &ResponseCounter{success: 1},
			},
		},
	}
	for i, c := range cases {
		got := c.base.Merge(c.input)
		if !reflect.DeepEqual(got, c.expect) {
			t.Errorf("#%d: want: %v, got: %v", i, c.expect, got)
		}
	}
}
