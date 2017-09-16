package server

import (
	"net/http"
	"net/url"
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
		res, err := client.PostForm(ts.URL+"/login", values)
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
	res, err := client.Get(ts.URL + "/login")
	if err != nil {
		t.Fatalf("want non error, got %v", err)
	}
	if res.StatusCode != http.StatusOK {
		t.Errorf("want %d, got %d", http.StatusOK, res.StatusCode)
	}
}
