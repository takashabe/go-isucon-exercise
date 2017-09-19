package server

import (
	"io/ioutil"
	"net/url"
	"strings"
	"testing"
)

func TestBenchDetail(t *testing.T) {
	ts := setupServer(t, "")
	defer ts.Close()
	setupFixture(t, "fixture/teams.yaml", "fixture/scores.yaml")

	values := url.Values{}
	values.Add("email", "foo")
	values.Add("password", "foo")
	client := clientWithNonRedirect()
	client.Jar = login(t, ts, values)

	cases := []struct {
		input        string
		expectPrefix []byte
	}{
		{"1", []byte(`{"id":1,"summary":"success","score":100,`)},
	}
	for i, c := range cases {
		res, err := client.Get(ts.URL + "/bench_detail/" + c.input)
		if err != nil {
			t.Fatalf("#%d: want not error, got %v", i, err)
		}
		defer res.Body.Close()

		payload, err := ioutil.ReadAll(res.Body)
		if err != nil {
			t.Fatalf("#%d: want not error, got %v", i, err)
		}
		if !strings.HasPrefix(string(payload), string(c.expectPrefix)) {
			t.Errorf("#%d: want has prefix %s, got %s", i, c.expectPrefix, payload)
		}
	}
}

func TestHistory(t *testing.T) {
	ts := setupServer(t, "")
	defer ts.Close()
	setupFixture(t, "fixture/teams.yaml", "fixture/scores.yaml")

	values := url.Values{}
	values.Add("email", "foo")
	values.Add("password", "foo")
	client := clientWithNonRedirect()
	client.Jar = login(t, ts, values)

	cases := []struct {
		input        string
		expectPrefix []byte
	}{
		{"1", []byte(`{"id":1,"summary":"success","score":100,`)},
	}
	for i, c := range cases {
		res, err := client.Get(ts.URL + "/bench_detail/" + c.input)
		if err != nil {
			t.Fatalf("#%d: want not error, got %v", i, err)
		}
		defer res.Body.Close()

		payload, err := ioutil.ReadAll(res.Body)
		if err != nil {
			t.Fatalf("#%d: want not error, got %v", i, err)
		}
		if !strings.HasPrefix(string(payload), string(c.expectPrefix)) {
			t.Errorf("#%d: want has prefix %s, got %s", i, c.expectPrefix, payload)
		}
	}
}
