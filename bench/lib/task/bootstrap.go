package task

import (
	"fmt"

	"github.com/PuerkitoBio/goquery"
)

// BootstrapTask checks initial content consistency
type BootstrapTask struct {
	w Worker
}

func (t *BootstrapTask) FinishHook(r Result) Result {
	r.Valid = true

	if len(r.Violations) > 0 {
		r.Valid = false
	}
	return r
}

func (t *BootstrapTask) Task(sessions []*Session) {
	// BootstrapTask use 0..2
	s1 := sessions[0]
	s2 := sessions[1]
	s3 := sessions[2]

	t.login2ndUser(s1)
	t.login3rdUser(s3)
	t.loginToIndex(s1)
	t.loginForm(s1)
	t.login1stUser(s1)
	t.indexAfterLoginDetail(s1)
	t.stylesheet(s1)
	t.indexAfterLogin(s2, "INDEX AFTER LOGIN 2ND USER")
	t.indexAfterLogin(s3, "INDEX AFTER LOGIN 3RD USER")
	t.postTweet(s1)
	t.viewProfileFollowUser(s1, s3)
	t.viewProfileNoFollowUser(s1, s2)
	t.postFollow(s2, s1)
	t.postTweetFromFollower(s1, s2)
	t.existFollower(s1, s2)
	t.logout(s1)
}

func (t *BootstrapTask) login2ndUser(s *Session) {
	t.w.getAndCheck(s, "/login", "LOGIN GET 2ND USER", func(c *Checker) {
		c.isStatusCode(200)
	})
	p := util.makeLoginParam(s.param.Email, s.param.Password)
	t.w.postAndCheck(s, "/login", p, "LOGIN POST 2ND USER", func(c *Checker) {
		c.isRedirect("/")
	})
}

func (t *BootstrapTask) login3rdUser(s *Session) {
	t.w.getAndCheck(s, "/login", "LOGIN GET 3RD USER", func(c *Checker) {
		c.isStatusCode(200)
	})
	p := util.makeLoginParam(s.param.Email, s.param.Password)
	t.w.postAndCheck(s, "/login", p, "LOGIN POST 3RD USER", func(c *Checker) {
		c.isRedirect("/")
	})
}

func (t *BootstrapTask) loginToIndex(s *Session) {
	t.w.getAndCheck(s, "/", "SHOULD LOGIN AT FIRST", func(c *Checker) {
		c.isRedirect("/login")
	})
}

func (t *BootstrapTask) loginForm(s *Session) {
	t.w.getAndCheck(s, "/", "LOGIN PAGE", func(c *Checker) {
		c.isStatusCode(200)

		c.nodeCount("form input[type=text]", 1)
		c.nodeCount("form input[type=password]", 1)
		c.nodeCount("form input[type=submit]", 1)
	})
}

func (t *BootstrapTask) login1stUser(s *Session) {
	p := util.makeLoginParam(s.param.Email, s.param.Password)
	t.w.postAndCheck(s, "/login", p, "LOGIN POST", func(c *Checker) {
		c.isRedirect("/")
	})
}

func (t *BootstrapTask) indexAfterLoginDetail(s *Session) {
	t.w.getAndCheck(s, "/", "INDEX AFTER LOGIN", func(c *Checker) {
		c.isStatusCode(200)

		c.hasStyleSheet("/css/bootstrap.min.css")

		c.hasContent("dd#prof-name", s.param.Name)
		c.hasContent("dd#prof-email", s.param.Email)

		c.nodeCount("dd#prof-following a", 1)
		c.attribute("dd#prof-following a", "href", "/following")
		c.nodeCount("dd#prof-followers a", 1)
		c.attribute("dd#prof-followers a", "href", "/followers")

		c.matchContent("dd#prof-followers", `\d`)
	})
}

func (t *BootstrapTask) stylesheet(s *Session) {
	t.w.getAndCheck(s, "/css/bootstrap.min.css", "STYLE SHEET CHECK", func(c *Checker) {
		c.isStatusCode(200)
		c.isContentLength(122540)
	})
}

func (t *BootstrapTask) indexAfterLogin(s *Session, requestName string) {
	t.w.getAndCheck(s, "/", requestName, func(c *Checker) {
		c.isStatusCode(200)
		c.hasStyleSheet("/css/bootstrap.min.css")
		c.hasContent("dd#prof-name", s.param.Name)
		c.hasNode("dd#prof-email")
	})
}

func (t *BootstrapTask) postTweet(s *Session) {
	p := util.makeTweetParam()
	t.w.postAndCheck(s, "/tweet", p, "POST NEW TWEET", func(c *Checker) {
		c.isRedirect("/")
	})
}

func (t *BootstrapTask) viewProfileFollowUser(s, dst *Session) {
	url := fmt.Sprintf("/user/%d", dst.param.ID)
	t.w.getAndCheck(s, url, "PROFILE FROM FOLLOW USER", func(c *Checker) {
		c.hasContent("dd#prof-name", s.param.Name)
		c.hasContent("dd#prof-email", s.param.Email)
		c.missingNode("form#follow-form")
	})
}

func (t *BootstrapTask) viewProfileNoFollowUser(s, dst *Session) {
	url := fmt.Sprintf("/user/%d", dst.param.ID)
	t.w.getAndCheck(s, url, "PROFILE FROM NON FOLLOW USER", func(c *Checker) {
		c.hasContent("dd#prof-name", s.param.Name)
		c.hasContent("dd#prof-email", s.param.Email)
		c.hasNode("form#follow-form")
	})
}

func (t *BootstrapTask) postFollow(s, dst *Session) {
	url := fmt.Sprintf("/follow/%d", dst.param.ID)
	p := util.makeTweetParam()
	t.w.postAndCheck(s, url, p, "POST FOLLOW", func(c *Checker) {
		c.isRedirect("/")
	})

	url = fmt.Sprintf("/user/%d", dst.param.ID)
	t.w.getAndCheck(s, url, "SEE 2ND USER FOLLOWING PAGE AFTER FOLLOW 1ST USER", func(c *Checker) {
		c.isStatusCode(200)
		c.contentFunc(
			fmt.Sprintf("#following dl dd.follow-follow a[href=%s]", url),
			"フォローしたばかりのユーザが含まれていません",
			func(se *goquery.Selection) bool {
				text, ok := se.Attr("href")
				return ok && text == url
			})
	})
}

func (t *BootstrapTask) postTweetFromFollower(s, dst *Session) {
	p := util.makeTweetParam()
	t.w.postAndCheck(s, "/tweet", p, "POST NEW TWEET", func(c *Checker) {
		c.isRedirect("/")
	})

	t.w.getAndCheck(dst, "/", "SEE 2ND USER TIMELINE AFTER TWEET 1ST USER", func(c *Checker) {
		c.isStatusCode(200)
		c.contentFunc(
			"#timeline.row.panel.panel-primary div.tweet div.tweet",
			"フォローしているユーザのツイートが含まれていません",
			func(se *goquery.Selection) bool {
				return se.Text() == p.Get("content")
			})
	})
}

func (t *BootstrapTask) existFollower(s, dst *Session) {
	t.w.getAndCheck(s, "/followers", "SEE 1ST USER FOLLOWERS PAGE AFTER FOLLOW FROM 2ND USER", func(c *Checker) {
		c.isStatusCode(200)
		c.contentFunc(
			"#followers.row.panel.panel-primary dl dd.follow-follow",
			"フォローされているユーザが含まれていません",
			func(se *goquery.Selection) bool {
				return se.Text() == s.param.Name
			})
	})
}

func (t *BootstrapTask) logout(s *Session) {
	t.w.getAndCheck(s, "/logout", "LOGOUT 1ST USER", func(c *Checker) {
		c.isRedirect("/login")
	})
	t.w.getAndCheck(s, "/", "INDEX AFTER LOGOUT", func(c *Checker) {
		c.isRedirect("/login")
	})
}
