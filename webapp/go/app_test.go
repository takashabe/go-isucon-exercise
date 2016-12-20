package main

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/k0kubun/pp"
)

func mockServer(handler func(http.ResponseWriter, *http.Request)) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(handler))
}

func TestLoginWithGet(t *testing.T) {
	server := mockServer(loginHandler)
	defer server.Close()

	res, err := http.Get(server.URL)
	if err != nil {
		t.Errorf("want no error, but %v", err.Error())
	}

	expectedCode := 200
	if res.StatusCode != expectedCode {
		t.Errorf("want %d, but %d", expectedCode, res.StatusCode)
	}

	if len(res.Cookies()) != 0 {
		t.Errorf("wont len 0, but len ", len(res.Cookies()))
	}
}

func TestLoginWithPost(t *testing.T) {
	server := mockServer(loginHandler)
	defer server.Close()

	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	// Success Authenticate
	auth := url.Values{}
	auth.Add("email", "Doris@example.com")
	auth.Add("password", "Doris")

	res, err := client.PostForm(server.URL, auth)
	if err != nil {
		t.Errorf("want no error, but %v", err.Error())
	}
	if res.StatusCode != 302 {
		t.Errorf("wont 302, but %d", res.StatusCode)
	}
	if loc, _ := res.Location(); loc.Path != "/" {
		t.Errorf("wont /, but %v", loc.Path)
	}
}

func TestLoginWithPost2(t *testing.T) {
	server := mockServer(loginHandler)
	defer server.Close()

	// Fail Authenticate
	// auth := url.Values{}
	res, err := http.Get(server.URL)
	// _, err := http.PostForm(server.URL, auth)
	if err != nil {
		t.Errorf("want no error, but %v", err.Error())
	}
	pp.Println(res)
	// if res.StatusCode != 401 {
	//   t.Errorf("wont 401, but %d", res.StatusCode)
	// }
}
