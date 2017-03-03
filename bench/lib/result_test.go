package bench

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
				RequestCount: 1,
				Response:     &ResponseCounter{success: 1},
				Violations:   make([]*Violation, 0),
			},
		},
		{
			302,
			&Result{
				RequestCount: 1,
				Response:     &ResponseCounter{redirect: 1},
				Violations:   make([]*Violation, 0),
			},
		},
		{
			404,
			&Result{
				RequestCount: 1,
				Response:     &ResponseCounter{clientError: 1},
				Violations:   make([]*Violation, 0),
			},
		},
		{
			999,
			&Result{
				RequestCount: 1,
				Response:     &ResponseCounter{serverError: 1},
				Violations:   make([]*Violation, 0),
			},
		},
		{
			0,
			&Result{
				RequestCount: 1,
				Response:     &ResponseCounter{serverError: 1},
				Violations:   make([]*Violation, 0),
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
				RequestCount: 1,
				Response:     &ResponseCounter{success: 200},
			},
			&Result{
				RequestCount: 2,
				Response:     &ResponseCounter{success: 200, exception: 1},
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
				Violations: []*Violation{&Violation{
					RequestName: "foo",
					Cause:       "bar",
					Count:       1,
				}},
			},
		},
		{
			"foo",
			"bar",
			&Result{
				Violations: []*Violation{
					&Violation{
						RequestName: "foo",
						Cause:       "bar",
						Count:       1,
					},
					&Violation{
						RequestName: "hoge",
						Cause:       "piyo",
						Count:       2,
					},
				},
			},
			&Result{
				Violations: []*Violation{
					&Violation{
						RequestName: "foo",
						Cause:       "bar",
						Count:       2,
					},
					&Violation{
						RequestName: "hoge",
						Cause:       "piyo",
						Count:       2,
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
				Valid:        false,
				RequestCount: 1,
				ElapsedTime:  300,
				Response:     &ResponseCounter{success: 1, exception: 1},
				Violations: []*Violation{
					&Violation{
						RequestName: "foo",
						Cause:       "bar",
						Count:       1,
					},
				},
			},
			Result{
				Valid:        true,
				RequestCount: 1,
				ElapsedTime:  300,
				Response:     &ResponseCounter{success: 1},
				Violations: []*Violation{
					&Violation{
						RequestName: "foo",
						Cause:       "bar",
						Count:       1,
					},
					&Violation{
						RequestName: "hoge",
						Cause:       "piyo",
						Count:       1,
					},
				},
			},
			&Result{
				Valid:        false,
				RequestCount: 2,
				ElapsedTime:  600,
				Response:     &ResponseCounter{success: 2, exception: 1},
				Violations: []*Violation{
					&Violation{
						RequestName: "foo",
						Cause:       "bar",
						Count:       2,
					},
					&Violation{
						RequestName: "hoge",
						Cause:       "piyo",
						Count:       1,
					},
				},
			},
		},
		{
			&Result{
				Response: &ResponseCounter{success: 1},
			},
			Result{
				Response: newResponse(),
			},
			&Result{
				Response: &ResponseCounter{success: 1},
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

func TestToJson(t *testing.T) {
	base := Result{
		Valid:        true,
		RequestCount: 10,
		ElapsedTime:  300,
		Response: &ResponseCounter{
			success:   5,
			exception: 5,
		},
		Violations: []*Violation{
			&Violation{
				RequestName: "a",
				Cause:       "b",
				Count:       2,
			},
			&Violation{
				RequestName: "c",
				Cause:       "d",
				Count:       2,
			},
		},
	}
	want := []byte(`{
	"valid": true,
	"request_count": 10,
	"elapsed_time": 300,
	"response": {},
	"violations": [
		{
			"request_type": "a",
			"description": "b",
			"num": 2
		},
		{
			"request_type": "c",
			"description": "d",
			"num": 2
		}
	]
}`)
	got, err := base.json()
	if err != nil {
		t.Errorf("want no error, got %v", err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("want %v, got %v", string(want), string(got))
	}
}
