package main

import (
	"database/sql"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"net/url"
	"strconv"
	"testing"

	"github.com/Puerkitobio/goquery"
	"github.com/takashabe/go-fixture"
	_ "github.com/takashabe/go-fixture/mysql"
	"github.com/takashabe/go-router"
)

func newRouter() http.Handler {
	r := router.NewRouter()
	r.Get("/", getIndex)
	r.Get("/login", getLogin)
	r.Get("/logout", getLogout)
	r.Get("/tweet", getTweet)
	r.Get("/user/:id", userHandler)
	r.Get("/following", getFollowing)
	r.Get("/followers", getFollowers)

	r.Post("/login", postLogin)
	r.Post("/tweet", postTweet)
	r.Post("/follow/:id", postFollow)

	return r
}

func LoadFixture(t *testing.T) {
	db, err := sql.Open("mysql", "isucon@/isucon_test")
	if err != nil {
		t.Fatal("failed to db open")
	}
	defer db.Close()

	f := fixture.NewFixture(db, "mysql")
	f.LoadSQL("fixture/create.sql")
	f.Load("fixture/user.yml")
	f.Load("fixture/tweet.yml")
	f.Load("fixture/follow.yml")
}

func getDummyLoginParams() url.Values {
	return url.Values{
		"email":    {"Helen@example.com"},
		"password": {"Helen"},
	}
}

func notRedirectClient() *http.Client {
	return &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
}

func cookieClient() *http.Client {
	jar, _ := cookiejar.New(nil)
	return &http.Client{
		Jar: jar,
	}
}

func login(t *testing.T, ts *httptest.Server) (*http.Client, *http.Response) {
	client := cookieClient()
	resp, err := client.PostForm(ts.URL+"/login", getDummyLoginParams())
	if err != nil {
		t.Fatalf("failed to login... err:%v", err)
	}
	return client, resp
}

func TestLoginGet(t *testing.T) {
	ts := httptest.NewServer(newRouter())
	defer ts.Close()

	resp, err := http.Get(ts.URL + "/login")
	if err != nil {
		t.Errorf("want: no error, got: %v", err.Error())
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		t.Errorf("want: %d, got: %d", 200, resp.StatusCode)
	}
}

func TestLoginPost(t *testing.T) {
	ts := httptest.NewServer(newRouter())
	defer ts.Close()

	cases := []struct {
		input            url.Values
		expectStatusCode int
		expectLocation   string
	}{
		{getDummyLoginParams(), 302, "/"},
		{url.Values{}, 401, ""},
		{url.Values{"email": {""}, "password": {""}}, 401, ""},
	}
	for i, c := range cases {
		client := notRedirectClient()
		resp, err := client.PostForm(ts.URL+"/login", c.input)
		if err != nil {
			t.Errorf("#%d: want no error, got %v", i, err.Error())
		}
		if c.expectStatusCode != resp.StatusCode {
			t.Errorf("#%d: want %d, got %d", i, c.expectStatusCode, resp.StatusCode)
		}
		if c.expectLocation != "" {
			if loc, _ := resp.Location(); c.expectLocation != loc.Path {
				t.Errorf("#%d: want %s, got %s", i, c.expectLocation, loc.Path)
			}
		}
	}
}

func TestIndexWithNotLogin(t *testing.T) {
	ts := httptest.NewServer(newRouter())
	defer ts.Close()

	client := notRedirectClient()
	res, err := client.Get(ts.URL)
	if err != nil {
		t.Errorf("want no error, but %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != 302 {
		t.Errorf("want 302, but %d", res.StatusCode)
	}
	if loc, err := res.Location(); err == nil && loc.Path != "/login" {
		t.Errorf("want /login, but %s", loc.Path)
	}
}

func TestIndexWithLogin(t *testing.T) {
	ts := httptest.NewServer(newRouter())
	defer ts.Close()

	client := cookieClient()

	// login
	loginResp, err := client.PostForm(ts.URL+"/login", getDummyLoginParams())
	if err != nil {
		t.Errorf("want no error, got %v", err)
	}
	defer loginResp.Body.Close()

	// get index after login
	indexResp, err := client.Get(ts.URL)
	if err != nil {
		t.Errorf("want no error, got %v", err)
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

func TestLoginAndLogout(t *testing.T) {
	ts := httptest.NewServer(newRouter())

	client, loginResp := login(t, ts)
	defer loginResp.Body.Close()

	indexResp, err := client.Get(ts.URL)
	if err != nil {
		t.Errorf("want no error, but %v", err)
	}
	if indexResp.StatusCode != 200 {
		t.Errorf("want 200, got %d", indexResp.StatusCode)
	}
	defer indexResp.Body.Close()

	// logout
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
	ts := httptest.NewServer(newRouter())
	defer ts.Close()

	// only redirect test
	client := notRedirectClient()

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
	ts := httptest.NewServer(newRouter())
	defer ts.Close()

	client, loginResp := login(t, ts)
	defer loginResp.Body.Close()

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
	ts := httptest.NewServer(newRouter())
	defer ts.Close()

	client := notRedirectClient()
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
	ts := httptest.NewServer(newRouter())
	defer ts.Close()

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
	ts := httptest.NewServer(newRouter())
	defer ts.Close()

	client := notRedirectClient()
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

func TestFollowingWithLogin(t *testing.T) {
	ts := httptest.NewServer(newRouter())
	defer ts.Close()

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

	resp, err := client.Get(ts.URL + "/following")
	defer resp.Body.Close()
	if err != nil {
		t.Errorf("want no error, got %v", err)
	}

	// parsed html body
	doc, err := goquery.NewDocumentFromResponse(resp)
	if err != nil {
		t.Errorf("want no error, got %v", err)
	}

	date := doc.Find("dt[class='follow-date']").Text()
	if len(date) <= 0 {
		t.Errorf("want len more than 0, got %s", date)
	}
}

func TestUserWithNotLogin(t *testing.T) {
	ts := httptest.NewServer(newRouter())
	defer ts.Close()

	client := notRedirectClient()
	resp, err := client.Get(ts.URL + "/user/30")
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

func TestUserWithLogin(t *testing.T) {
	cases := []struct {
		input             string
		expectStatusCode  int
		expectFollowExist bool
		expectTweetExist  bool
	}{
		{"/user/30", 200, false, true},
		{"/user/1000", 200, true, true},
		{"/user/100", 200, false, true},
		{"/user/0", 404, false, false},
	}
	for i, c := range cases {
		// login
		ts := httptest.NewServer(newRouter())
		defer ts.Close()
		jar, _ := cookiejar.New(nil)
		client := &http.Client{
			Jar:           jar,
			CheckRedirect: func(req *http.Request, via []*http.Request) error { return http.ErrUseLastResponse },
		}
		loginResp, err := client.PostForm(ts.URL+"/login", getDummyLoginParams())
		defer loginResp.Body.Close()
		if err != nil {
			t.Errorf("#%d want no error, got %v", i, err)
		}

		// test
		resp, err := client.Get(ts.URL + c.input)
		if err != nil {
			t.Errorf("#%d want no error, got %v", i, err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != c.expectStatusCode {
			t.Errorf("#%d want %d, got %d", i, c.expectStatusCode, resp.StatusCode)
		}
		doc, err := goquery.NewDocumentFromResponse(resp)
		if err != nil {
			t.Errorf("#%d want no error, got %v", i, err)
		}
		if v := doc.Find("form[id='follow-form']").Text(); c.expectFollowExist && len(v) <= 0 {
			t.Errorf("#%d want exist id='follow-form'", i)
		}
		if v := doc.Find("div[class='user']").Text(); c.expectTweetExist && len(v) <= 0 {
			t.Errorf("#%d want exist id='follow-form'", i)
		}
	}
}

func TestFollowersWithNotLogin(t *testing.T) {
	ts := httptest.NewServer(newRouter())
	defer ts.Close()

	client := notRedirectClient()
	resp, err := client.Get(ts.URL + "/user/30")
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

func TestFollowersWithLogin(t *testing.T) {
	ts := httptest.NewServer(newRouter())
	defer ts.Close()

	client, loginResp := login(t, ts)
	defer loginResp.Body.Close()

	resp, err := client.Get(ts.URL + "/followers")
	defer resp.Body.Close()
	if err != nil {
		t.Errorf("want no error, got %v", err)
	}
	doc, err := goquery.NewDocumentFromResponse(resp)
	if err != nil {
		t.Errorf("want no error, got %v", err)
	}
	date := doc.Find("dt[class='follow-date']").Text()
	if len(date) == 0 {
		t.Errorf("want len more than 0, got %s", date)
	}
}

func TestFollow(t *testing.T) {
	ts := httptest.NewServer(newRouter())
	defer ts.Close()

	client, loginResp := login(t, ts)
	client.CheckRedirect = func(req *http.Request, via []*http.Request) error { return http.ErrUseLastResponse }
	defer loginResp.Body.Close()

	req, err := http.NewRequest("POST", ts.URL+"/follow/100", nil)
	if err != nil {
		t.Errorf("want no error, got %v", err)
	}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("want no error, got %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 303 {
		t.Errorf("want 303, got %d", resp.StatusCode)
	}
	if loc, err := resp.Location(); err == nil && loc.Path != "/" {
		t.Errorf("want /, got %s", loc.Path)
	}
}
