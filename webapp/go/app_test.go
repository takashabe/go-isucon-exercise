package main

import (
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"net/url"
	"strconv"
	"testing"

	"github.com/Puerkitobio/goquery"
)

func TestLoginWithGet(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(loginHandler))
	defer server.Close()

	res, err := http.Get(server.URL)
	if err != nil {
		t.Errorf("want no error, but %v", err)
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
	ts := httptest.NewServer(http.HandlerFunc(loginHandler))
	defer ts.Close()

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
	res, err := client.PostForm(ts.URL, auth)
	if err != nil {
		t.Errorf("want no error, but %v", err)
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
	res, err = client.PostForm(ts.URL, emptyAuth)
	if err != nil {
		t.Errorf("want no error, but %v", err)
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
	res, err = client.PostForm(ts.URL, invalidAuth)
	if err != nil {
		t.Errorf("want no error,  but %v", err)
	}
	if res.StatusCode != 401 {
		t.Errorf("invalidAuth: want 401,  but %d", res.StatusCode)
	}
	defer res.Body.Close()
}

func TestIndexWithNotLogin(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(indexHandler))
	defer ts.Close()

	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	res, err := client.Get(ts.URL)
	if err != nil {
		t.Errorf("want no error, but %v", err)
	}
	if res.StatusCode != 302 {
		t.Errorf("want 302, but %d", res.StatusCode)
	}
	if loc, _ := res.Location(); loc.Path != "/login" {
		t.Errorf("want /login, but %s", loc.Path)
	}
	defer res.Body.Close()
}

func TestIndexWithLogin(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/", indexHandler)
	mux.HandleFunc("/login", loginHandler)

	ts := httptest.NewServer(mux)
	defer ts.Close()

	jar, _ := cookiejar.New(nil)
	client := &http.Client{
		Jar: jar,
	}

	// login
	auth := url.Values{
		"email":    {"Doris@example.com"},
		"password": {"Doris"},
	}
	loginResp, err := client.PostForm(ts.URL+"/login", auth)
	if err != nil {
		t.Errorf("want no error, but %v", err)
	}
	defer loginResp.Body.Close()

	// get index after login
	indexResp, err := client.Get(ts.URL)
	if err != nil {
		t.Errorf("want no error, but %v", err)
	}
	defer indexResp.Body.Close()

	// parsed html body
	doc, err := goquery.NewDocumentFromResponse(indexResp)
	if err != nil {
		t.Errorf("want no error, got %v", err)
	}

	name := doc.Find("dd[id='prof-name']").Text()
	wantName := "Doris"
	if name != wantName {
		t.Errorf("want %s, got %s", wantName, name)
	}

	follow := doc.Find("dd[id='prof-following']").Text()
	if i, _ := strconv.Atoi(follow); i < 0 {
		t.Errorf("want follow count more than 0, got %s", follow)
	}

	followers := doc.Find("dd[id='prof-followers']").Text()
	if i, _ := strconv.Atoi(followers); i <= 0 {
		t.Errorf("want followers count more than 0,  got %s", followers)
	}

	tweet := doc.Find("div[class='tweet']").First().Text()
	if len(tweet) <= 0 {
		t.Errorf("want len more than 0, got %s", tweet)
	}
}

func TestLogoutHandler(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/", indexHandler)
	mux.HandleFunc("/login", loginHandler)
	mux.HandleFunc("/logout", logoutHandler)

	ts := httptest.NewServer(mux)

	jar, _ := cookiejar.New(nil)
	client := &http.Client{
		Jar: jar,
	}

	auth := url.Values{
		"email":    {"Doris@example.com"},
		"password": {"Doris"},
	}
	loginResp, err := client.PostForm(ts.URL+"/login", auth)
	if err != nil {
		t.Errorf("want no error, but %v", err)
	}
	defer loginResp.Body.Close()

	indexResp, err := client.Get(ts.URL)
	if err != nil {
		t.Errorf("want no error, but %v", err)
	}
	if indexResp.StatusCode != 200 {
		t.Errorf("want 200, got %d", indexResp.StatusCode)
	}
	defer indexResp.Body.Close()

	// test http headers on redirect
	client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}
	logoutResp, err := client.Get(ts.URL + "/logout")
	if logoutResp.StatusCode != 302 {
		t.Errorf("want 302, got %d", logoutResp.StatusCode)
	}
	if loc, _ := logoutResp.Location(); loc.Path != "/login" {
		t.Errorf("want /login, got %s", loc.Path)
	}
	defer logoutResp.Body.Close()

	indexResp2, err := client.Get(ts.URL)
	if indexResp2.StatusCode != 302 {
		t.Errorf("want 302, got %d", indexResp2.StatusCode)
	}
	if loc, _ := indexResp2.Location(); loc.Path != "/login" {
		t.Errorf("want /login, got %s", loc.Path)
	}
	defer indexResp2.Body.Close()
}

func TestTweetWithGet(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(tweetHandler))
	defer ts.Close()

	resp, err := http.Get(ts.URL)
	defer resp.Body.Close()
	if err != nil {
		t.Errorf("want no error, got %v", err)
	}

	if resp.StatusCode != 200 {
		t.Errorf("want 200, got %d", resp.StatusCode)
	}
}
