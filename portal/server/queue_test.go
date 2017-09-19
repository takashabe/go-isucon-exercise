package server

import (
	"encoding/json"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/takashabe/go-isucon-exercise/portal/models"
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

func dummyEnqueueWithTeam(t *testing.T, ts *httptest.Server, email, password string) {
	values := url.Values{}
	values.Add("email", email)
	values.Add("password", password)
	client := clientWithNonRedirect()
	client.Jar = login(t, ts, values)
	res, err := client.Post(ts.URL+"/enqueue", "application/json", nil)
	if err != nil {
		t.Fatalf("want non error, got %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Fatalf("want status code %d, got %d", http.StatusOK, res.StatusCode)
	}
}

func TestQueues(t *testing.T) {
	pubsubServer := setupPubsub(t)
	defer pubsubServer.Close()
	portalServer := setupServer(t, pubsubServer.URL)
	defer portalServer.Close()
	setupFixture(t, "fixture/teams.yaml")
	setupFixture(t, "fixture/queues.yaml")
	dummyEnqueueWithTeam(t, portalServer, "foo", "foo")
	dummyEnqueueWithTeam(t, portalServer, "bar", "bar")

	values := url.Values{}
	values.Add("email", "foo")
	values.Add("password", "foo")
	client := clientWithNonRedirect()
	client.Jar = login(t, portalServer, values)
	res, err := client.Get(portalServer.URL + "/queues")
	if err != nil {
		t.Fatalf("want non error, got %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Fatalf("want status code %d, got %d", http.StatusOK, res.StatusCode)
	}
	var decorder []models.CurrentQueue
	err = json.NewDecoder(res.Body).Decode(&decorder)
	if err != nil {
		t.Fatalf("want non error, got %v", err)
	}

	var exist bool
	for _, d := range decorder {
		if d.MyTeam {
			exist = true
			break
		}
	}
	if !exist {
		t.Errorf("want contain MyTeam sent queue, got %v", decorder)
	}
}
