package main

import (
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
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
		t.Errorf("auth: want 302, but %d", res.StatusCode)
	}
	if loc, _ := res.Location(); loc.Path != "/" {
		t.Errorf("want /, but %v", loc.Path)
	}
	res.Body.Close()

	// Empty login params
	emptyAuth := url.Values{}
	res, err = client.PostForm(server.URL, emptyAuth)
	if err != nil {
		t.Errorf("want no error, but %v", err.Error())
	}
	if res.StatusCode != 401 {
		t.Errorf("emptyAuth: want 401, but %d", res.StatusCode)
	}
	res.Body.Close()

	// Invalid login params
	invalidAuth := url.Values{
		"email":    {"empty"},
		"password": {"empty"},
	}
	res, err = client.PostForm(server.URL, invalidAuth)
	if err != nil {
		t.Errorf("want no error,  but %v", err.Error())
	}
	if res.StatusCode != 401 {
		t.Errorf("invalidAuth: want 401,  but %d", res.StatusCode)
	}
	res.Body.Close()
}

func TestIndex(t *testing.T) {
	// when non login
	server := httptest.NewServer(http.HandlerFunc(indexHandler))
	defer server.Close()

	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	res, err := client.Get(server.URL)
	if err != nil {
		t.Errorf("want no error, but %v", err.Error())
	}
	if res.StatusCode != 302 {
		t.Errorf("want 302, but %d", res.StatusCode)
	}
	if loc, _ := res.Location(); loc.Path != "/login" {
		t.Errorf("want /login, but %s", loc.Path)
	}
	res.Body.Close()

	// login to index
	loginServer := httptest.NewServer(http.HandlerFunc(loginHandler))
	defer loginServer.Close()
	indexServer := httptest.NewServer(http.HandlerFunc(indexHandler))
	defer indexServer.Close()

	jar, _ := cookiejar.New(nil)
	client = &http.Client{
		Jar: jar,
	}

	auth := url.Values{
		"email":    {"Doris@example.com"},
		"password": {"Doris"},
	}
	res, err = client.PostForm(loginServer.URL, auth)
	if err != nil {
		t.Errorf("want no error, but %v", err.Error())
	}
	res.Body.Close()

	res, err = client.Get(indexServer.URL)
	if err != nil {
		t.Errorf("want no error, but %v", err.Error())
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Error("want no error, but %v", err.Error())
	}

	// TODO: bodyパースしてツイート情報のエレメントがあるかどうかをテストする

	res.Body.Close()
}
