package server

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"reflect"
	"testing"
)

func TestLogin(t *testing.T) {
	setupFixture(t, "fixture/schema.sql", "fixture/teams.yaml")
	ts := setupServer(t, "")
	defer ts.Close()

	cases := []struct {
		email      string
		password   string
		location   string
		expectCode int
	}{
		{
			"foo",
			"foo",
			ts.URL + "/",
			http.StatusFound,
		},
		{
			"foo",
			"bar",
			"",
			http.StatusUnauthorized,
		},
	}
	for i, c := range cases {
		values := url.Values{}
		values.Add("email", c.email)
		values.Add("password", c.password)

		client := clientWithNonRedirect()
		res, err := client.PostForm(ts.URL+"/api/login", values)
		if err != nil {
			t.Fatalf("#%d: want non error, got %v", i, err)
		}
		defer res.Body.Close()

		if c.expectCode != res.StatusCode {
			t.Errorf("#%d: want %d, got %d", i, c.expectCode, res.StatusCode)
		}

		if c.expectCode == http.StatusFound {
			l, err := res.Location()
			if err != nil {
				t.Fatalf("#%d: want non error, got %v", i, err)
			}
			if c.location != l.String() {
				t.Errorf("#%d: want %s, got %s", i, c.location, l)
			}
		}
	}
}

func TestLogout(t *testing.T) {
	ts := setupServer(t, "")
	defer ts.Close()

	client := clientWithNonRedirect()
	res, err := client.Get(ts.URL + "/api/login")
	if err != nil {
		t.Fatalf("want non error, got %v", err)
	}
	if res.StatusCode != http.StatusOK {
		t.Errorf("want %d, got %d", http.StatusOK, res.StatusCode)
	}
}

func TestGetTeam(t *testing.T) {
	ts := setupServer(t, "")
	defer ts.Close()

	client := clientWithNonRedirect()
	values := url.Values{}
	values.Add("email", "foo")
	values.Add("password", "foo")
	client.Jar = login(t, ts, values)

	res, err := client.Get(ts.URL + "/api/team")
	if err != nil {
		t.Fatalf("want non error, got %v", err)
	}
	defer res.Body.Close()

	payload, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatalf("want non error, got %v", err)
	}
	expect := []byte(`{"ID":1,"Name":"team1","Instance":"localhost:8080"}`)
	if !reflect.DeepEqual(expect, payload) {
		t.Errorf("want %s, got %s", expect, payload)
	}
}
