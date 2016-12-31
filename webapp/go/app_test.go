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

func getServerMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", indexHandler)
	mux.HandleFunc("/login", loginHandler)
	mux.HandleFunc("/logout", logoutHandler)
	mux.HandleFunc("/tweet", tweetHandler)
	mux.HandleFunc("/following", followingHandler)
	return mux
}

func getDummyLoginParams() url.Values {
	auth := url.Values{
		"email":    {"Doris@example.com"},
		"password": {"Doris"},
	}
	return auth
}

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

	// only redirect test
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error { return http.ErrUseLastResponse },
	}

	// Success Authenticate
	res, err := client.PostForm(ts.URL, getDummyLoginParams())
	if err != nil {
		t.Errorf("want no error, but %v", err)
	}
	if res.StatusCode != 302 {
		t.Errorf("auth: want 302, but %d", res.StatusCode)
	}
	if loc, err := res.Location(); err == nil && loc.Path != "/" {
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
		"email":    {""},
		"password": {""},
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
	ts := httptest.NewServer(getServerMux())
	defer ts.Close()

	// only redirect test
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error { return http.ErrUseLastResponse },
	}

	res, err := client.Get(ts.URL)
	if err != nil {
		t.Errorf("want no error, but %v", err)
	}
	if res.StatusCode != 302 {
		t.Errorf("want 302, but %d", res.StatusCode)
	}
	if loc, err := res.Location(); err == nil && loc.Path != "/login" {
		t.Errorf("want /login, but %s", loc.Path)
	}
	defer res.Body.Close()
}

func TestIndexWithLogin(t *testing.T) {
	ts := httptest.NewServer(getServerMux())
	defer ts.Close()

	jar, _ := cookiejar.New(nil)
	client := &http.Client{
		Jar: jar,
	}

	// login
	loginResp, err := client.PostForm(ts.URL+"/login", getDummyLoginParams())
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
	ts := httptest.NewServer(getServerMux())

	jar, _ := cookiejar.New(nil)
	client := &http.Client{
		Jar: jar,
	}

	loginResp, err := client.PostForm(ts.URL+"/login", getDummyLoginParams())
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

	// only redirect test
	client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}
	logoutResp, err := client.Get(ts.URL + "/logout")
	if logoutResp.StatusCode != 302 {
		t.Errorf("want 302, got %d", logoutResp.StatusCode)
	}
	if loc, err := logoutResp.Location(); err == nil && loc.Path != "/login" {
		t.Errorf("want /login, got %s", loc.Path)
	}
	defer logoutResp.Body.Close()

	indexResp2, err := client.Get(ts.URL)
	if indexResp2.StatusCode != 302 {
		t.Errorf("want 302, got %d", indexResp2.StatusCode)
	}
	if loc, err := indexResp2.Location(); err == nil && loc.Path != "/login" {
		t.Errorf("want /login, got %s", loc.Path)
	}
	defer indexResp2.Body.Close()
}

func TestTweetWithNotLogin(t *testing.T) {
	ts := httptest.NewServer(getServerMux())
	defer ts.Close()

	// only redirect test
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error { return http.ErrUseLastResponse },
	}

	// GET
	getResp, err := client.Get(ts.URL + "/tweet")
	defer getResp.Body.Close()
	if err != nil {
		t.Errorf("want no error, got %v", err)
	}
	if getResp.StatusCode != 302 {
		t.Errorf("want 302, got %d", getResp.StatusCode)
	}
	if loc, err := getResp.Location(); err == nil && loc.Path != "/login" {
		t.Errorf("want /login, got %s", loc.Path)
	}

	// POST
	tweet := url.Values{
		"content": {"hello"},
	}
	postResp, err := client.PostForm(ts.URL+"/tweet", tweet)
	if err != nil {
		t.Errorf("want no error, got %v", err)
	}
	if postResp.StatusCode != 303 {
		t.Errorf("want 303, got %d", postResp.StatusCode)
	}
	if loc, err := postResp.Location(); err == nil && loc.Path != "/login" {
		t.Errorf("want /login, got %s", loc.Path)
	}
}

func TestTweetWithLoginGet(t *testing.T) {
	ts := httptest.NewServer(getServerMux())
	defer ts.Close()

	jar, _ := cookiejar.New(nil)
	client := &http.Client{
		Jar: jar,
	}

	loginResp, err := client.PostForm(ts.URL+"/login", getDummyLoginParams())
	defer loginResp.Body.Close()
	if err != nil {
		t.Errorf("want no error, got %v", err)
	}

	resp, err := client.Get(ts.URL + "/tweet")
	defer resp.Body.Close()
	if err != nil {
		t.Errorf("want no error, got %v", err)
	}
	if resp.StatusCode != 200 {
		t.Errorf("want 200, got %d", resp.StatusCode)
	}
}

func TestTweetWithNotLoginPost(t *testing.T) {
	ts := httptest.NewServer(getServerMux())
	defer ts.Close()

	// only redirect test
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error { return http.ErrUseLastResponse },
	}

	tweet := url.Values{
		"content": {"hello"},
	}
	tweetResp, err := client.PostForm(ts.URL+"/tweet", tweet)
	defer tweetResp.Body.Close()
	if err != nil {
		t.Errorf("want no error, got %v", err)
	}
	if tweetResp.StatusCode != 303 {
		t.Errorf("want 303, got %d", tweetResp.StatusCode)
	}
	if loc, err := tweetResp.Location(); err == nil && loc.Path != "/login" {
		t.Errorf("want /login, got %s", loc.Path)
	}
}

func TestTweetWithLoginPost(t *testing.T) {
	ts := httptest.NewServer(getServerMux())
	defer ts.Close()

	// only redirect test
	jar, _ := cookiejar.New(nil)
	client := &http.Client{
		Jar:           jar,
		CheckRedirect: func(req *http.Request, via []*http.Request) error { return http.ErrUseLastResponse },
	}

	loginResp, err := client.PostForm(ts.URL+"/login", getDummyLoginParams())
	defer loginResp.Body.Close()
	if err != nil {
		t.Errorf("want no error, got %v", err)
	}

	tweet := url.Values{
		"content": {"hello"},
	}
	tweetResp, err := client.PostForm(ts.URL+"/tweet", tweet)
	defer tweetResp.Body.Close()
	if err != nil {
		t.Errorf("want no error, got %v", err)
	}
	if tweetResp.StatusCode != 303 {
		t.Errorf("want 303, got %d", tweetResp.StatusCode)
	}
	if loc, err := tweetResp.Location(); err == nil && loc.Path != "/" {
		t.Errorf("want /, got %s", loc.Path)
	}
}

func TestFollowingWithNotLogin(t *testing.T) {
	ts := httptest.NewServer(getServerMux())
	defer ts.Close()

	// only redirect test
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error { return http.ErrUseLastResponse },
	}

	resp, err := client.Get(ts.URL + "/following")
	defer resp.Body.Close()
	if err != nil {
		t.Errorf("want no error, got %v", err)
	}
	if resp.StatusCode != 302 {
		t.Errorf("want 302, got %d", resp.StatusCode)
	}
	if loc, err := resp.Location(); err == nil && loc.Path != "/login" {
		t.Errorf("want /login, got %s", loc.Path)
	}
}
