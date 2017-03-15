package main

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"golang.org/x/net/html"

	"github.com/PuerkitoBio/goquery"
)

// Make load and check request. NOT SUPPORT CONCURRENCY
type LoadCheckerTask struct {
	timeout time.Time
}

func (t *LoadCheckerTask) FinishHook(r Result) Result {
	if len(r.Violations) > 0 {
		r.Fail()
	}
	return r
}

func (t *LoadCheckerTask) Task(ctx Ctx, d *Driver) *Driver {
	runningTime := ctx.workerRunningTime
	t.timeout = time.Now().Add(time.Millisecond * time.Duration(runningTime))
	for {
		if t.isTimeout() {
			return d
		}
		t.run(ctx, d)
	}
}

func (t *LoadCheckerTask) isTimeout() bool {
	return t.timeout.Before(time.Now())
}

func (t *LoadCheckerTask) run(ctx Ctx, d *Driver) {
	// 0..2 Bootstrap
	// 3..9 LoadChecker
	// 10.. Load
	rand.Seed(time.Now().UnixNano())
	sub := ctx.sessions[3:9]
	s1 := sub[rand.Intn(len(sub))]
	s2 := sub[rand.Intn(len(sub))]

	// login from s1
	d.get(s1, "/logout")
	d.post(s1, "/login", util.makeLoginParam(s1.param.Email, s1.param.Password))
	if t.isTimeout() {
		return
	}

	// check login s1
	if d.getAndStatus(s1, "/login") != 200 {
		d.getAndCheck(s1, "/login", "LOGIN PAGE BECAUSE NOT LOGGED IN", func(c *Checker) {
			c.isStatusCode(200)
		})
		d.postAndCheck(s1, "/login", util.makeLoginParam(s1.param.Email, s1.param.Password), "LOGIN POST WHEN LOGGED OUT", func(c *Checker) {
			c.isRedirect("/")
		})
		d.getAndCheck(s1, "/", "SHOW INDEX AFTER LOGIN", func(c *Checker) {
			c.isStatusCode(200)
		})
	}
	if t.isTimeout() {
		return
	}

	// check s2 long scenario
	// login from s2
	if d.getAndStatus(s2, "/") != 200 {
		d.getAndCheck(s2, "/login", "LOGIN PAGE BECAUSE NOT LOGGED IN", func(c *Checker) {
			c.isStatusCode(200)
		})
		d.postAndCheck(s2, "/login", util.makeLoginParam(s2.param.Email, s2.param.Password), "LOGIN POST WHEN LOGGED OUT", func(c *Checker) {
			c.isRedirect("/")
		})
		d.getAndCheck(s2, "/", "SHOW INDEX AFTER LOGIN", func(c *Checker) {
			c.isStatusCode(200)
		})
	}
	if t.isTimeout() {
		return
	}

	// use with s1
	followPath := fmt.Sprintf("/user/%d", s2.param.ID)
	name := d.getAndContent(s1, "/following",
		fmt.Sprintf("#following dl dd.follow-follow a[href='%s']", followPath), 0, func(node *html.Node) string {
			return node.Attr[0].Val
		})
	if name != "" {
		// already following, do tweet
		tweet := util.makeTweetParam()
		d.postAndCheck(s2, "/tweet", tweet, "POST TWEET FROM FOLLOWING USER", func(c *Checker) {
			c.isRedirect("/")
			// TODO: c.responseUntil
		})
		if t.isTimeout() {
			return
		}

		d.getAndCheck(s1, "/", "SEE FOLLOWING TWEET", func(c *Checker) {
			c.isStatusCode(200)
			c.contentFunc(
				"#timeline.row.panel.panel-primary div.tweet div.tweet",
				"フォローしているユーザのツイートが含まれていません",
				func(se *goquery.Selection) bool {
					text := strings.TrimSpace(se.Text())
					return text == tweet.Get("content")
				})
		})
	} else {
		// not yet following, do follow
		d.postAndCheck(s1, fmt.Sprintf("/follow/%d", s2.param.ID), nil, "MAKE FOLLOW", func(c *Checker) {
			c.isRedirect("/")
			// TODO: c.responseUntil
		})
		if t.isTimeout() {
			return
		}

		d.getAndCheck(s1, "following", "FOLLOWING LIST AFTER MAKING FOLLOW", func(c *Checker) {
			c.isStatusCode(200)
			c.contentFunc(
				fmt.Sprintf("#following dl dd.follow-follow a[href='%s']", followPath),
				"フォローしたばかりのユーザが含まれていません",
				func(se *goquery.Selection) bool {
					text, ok := se.Attr("href")
					return ok && text == followPath
				})
		})
	}
}
