package main

import (
	"log"
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

	doc, err := goquery.NewDocumentFromResponse(res)
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

	res.Body.Close()
}

func TestLogoutHandler(t *testing.T) {
	// login to index
	loginServer := httptest.NewServer(http.HandlerFunc(loginHandler))
	defer loginServer.Close()
	indexServer := httptest.NewServer(http.HandlerFunc(indexHandler))
	defer indexServer.Close()
	logoutServer := httptest.NewServer(http.HandlerFunc(logoutHandler))
	defer logoutServer.Close()

	jar, _ := cookiejar.New(nil)
	cookieClient := &http.Client{
		Jar: jar,
	}

	auth := url.Values{
		"email":    {"Doris@example.com"},
		"password": {"Doris"},
	}
	res, err := cookieClient.PostForm(loginServer.URL, auth)
	if err != nil {
		t.Errorf("want no error, but %v", err.Error())
	}
	res.Body.Close()
	log.Println("cookieClient.postform(login)")

	res, err = cookieClient.Get(indexServer.URL)
	if err != nil {
		t.Errorf("want no error, but %v", err.Error())
	}
	if res.StatusCode != 200 {
		t.Errorf("want 200, got %d", res.StatusCode)
	}
	res.Body.Close()
	log.Println("cookieClient.get(index)")

	res, err = cookieClient.Get(logoutServer.URL)
	if res.StatusCode != 302 {
		t.Errorf("want 302, got %d", res.StatusCode)
	}
	if loc, _ := res.Location(); loc.Path != "/login" {
		t.Errorf("want /login, got %s", loc.Path)
	}
	res.Body.Close()
	log.Println("cookieClient.get(logout)")

	cookieClient.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}
	res, err = cookieClient.Get(indexServer.URL)
	if res.StatusCode != 302 {
		t.Errorf("want 302, got %d", res.StatusCode)
	}
	if loc, _ := res.Location(); loc.Path != "/login" {
		t.Errorf("want /login, got %s", loc.Path)
	}
	res.Body.Close()
	log.Println("cookieClient.get(index)")
}
