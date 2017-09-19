package server

import (
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"net/url"
	"testing"
)

func login(t *testing.T, ts *httptest.Server, values url.Values) *cookiejar.Jar {
	jar, err := cookiejar.New(nil)
	if err != nil {
		t.Fatalf("want non error, got %v", err)
	}
	client := clientWithNonRedirect()
	client.Jar = jar
	res, err := client.PostForm(ts.URL+"/login", values)
	if err != nil {
		t.Fatalf("want non error, got %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusFound {
		t.Fatalf("want status code %d, got %d", http.StatusFound, res.StatusCode)
	}
	return jar
}

func TestEnqueue(t *testing.T) {
	pubsubServer := setupPubsub(t)
	defer pubsubServer.Close()
	portalServer := setupServer(t, pubsubServer.URL)
	defer portalServer.Close()
	setupFixture(t, "fixture/teams.yaml")

	values := url.Values{}
	values.Add("email", "foo")
	values.Add("password", "foo")
	jar := login(t, portalServer, values)

	client := clientWithNonRedirect()
	client.Jar = jar
	res, err := client.Post(portalServer.URL+"/enqueue", "application/json", nil)
	if err != nil {
		t.Fatalf("want non error, got %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Errorf("want status code %d, got %d", http.StatusOK, res.StatusCode)
	}
}
