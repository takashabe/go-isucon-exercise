package main

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestLoginWithGet(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(loginHandler))
	defer server.Close()

	res, err := http.Get(server.URL)
	if err != nil {
		t.Errorf("want no error, but %v", err.Error())
	}
	defer res.Body.Close()

	expectedCode := 200
	if res.StatusCode != expectedCode {
		t.Errorf("want %d, but %d", expectedCode, res.StatusCode)
	}
	if len(res.Cookies()) != 0 {
		t.Errorf("wont len 0, but len ", len(res.Cookies()))
	}
}

func TestLoginWithPost(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(loginHandler))
	defer server.Close()

	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	// Success Authenticate
	auth := url.Values{
		"email":    {"Doris@example.com"},
		"password": {"Doris"},
	}
	res, err := client.PostForm(server.URL, auth)
	if err != nil {
		t.Errorf("want no error, but %v", err.Error())
	}
	if res.StatusCode != 302 {
		t.Errorf("want 302, but %d", res.StatusCode)
	}
	if loc, _ := res.Location(); loc.Path != "/" {
		t.Errorf("want /, but %v", loc.Path)
	}
	res.Body.Close()

	// Failure Authenticate
	emptyAuth := url.Values{}
	res, err = client.PostForm(server.URL, emptyAuth)
	if err != nil {
		t.Errorf("want no error, but %v", err.Error())
	}
	if res.StatusCode != 401 {
		t.Errorf("want 401, but %d", res.StatusCode)
	}
	res.Body.Close()
}
