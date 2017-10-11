package benchmark

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"

	"golang.org/x/net/html"

	"github.com/pkg/errors"
	"github.com/takashabe/go-microbenchmark"
)

// request error violation errors
var (
	causeFailedReceiveResponse = "パス '%s' からレスポンスが返ってきませんでした"
)

type Driver struct {
	result *Result
	ctx    Ctx
}

func (d *Driver) get(sess *Session, path string) {
	d.getAndCheck(sess, path, "", nil)
}

func (d *Driver) getAndStatus(sess *Session, path string) int {
	var statusCode int
	d.getAndCheck(sess, path, "TO READ STATUS", func(c *Checker) {
		statusCode = c.statusCode()
	})
	return statusCode
}

func (d *Driver) getAndContent(sess *Session, path, selector string, i int, f func(node *html.Node) string) string {
	var content string
	d.getAndCheck(sess, path, "TO READ NODE", func(c *Checker) {
		doc, err := c.getDocument()
		if err != nil {
			return
		}
		sec := doc.Find(selector)
		if sec.Size() > i {
			content = f(sec.Get(i))
		}
	})
	return content
}

var tr = &http.Transport{
	MaxIdleConns:        10,
	MaxIdleConnsPerHost: 10,
}

func (d *Driver) getAndCheck(sess *Session, path, requestName string, check func(c *Checker)) {
	req, err := http.NewRequest("GET", d.ctx.uri(path), nil)
	if err != nil {
		log.Println(errors.Errorf("failed to generate request: %v", err.Error()))
		return
	}

	client := &http.Client{
		Transport: tr,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	if sess != nil {
		client.Jar = sess.cookie
	}
	d.requestAndCheck(path, requestName, req, client, check)
}

func (d *Driver) post(sess *Session, path string, params url.Values) {
	d.postAndCheck(sess, path, params, "", nil)
}

func (d *Driver) postAndCheck(sess *Session, path string, params url.Values, requestName string, check func(c *Checker)) {
	req, err := http.NewRequest("POST", d.ctx.uri(path), strings.NewReader(params.Encode()))
	if err != nil {
		log.Println(errors.Errorf("failed to generate request: %v", err.Error()))
		return
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{
		Transport: tr,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	if sess != nil {
		client.Jar = sess.cookie
	}
	d.requestAndCheck(path, requestName, req, client, check)
}

func (d *Driver) requestAndCheck(path, requestName string, req *http.Request, client *http.Client, check func(c *Checker)) {
	bench := microbenchmark.NewBenchmark()
	bench.Begin()
	res, err := client.Do(req)
	if err != nil {
		PrintDebugf("failed to response. path=%s, error=%v", path, err)
		d.result.addViolation(requestName, fmt.Sprintf(causeFailedReceiveResponse, path))
		d.result.addResponse(500)
		return
	}
	defer res.Body.Close()
	time := bench.End()

	d.result.addResponse(res.StatusCode)
	d.result.ElapsedTime = time.Nanoseconds()
	if check != nil {
		check(&Checker{
			ctx:          d.ctx,
			result:       d.result,
			path:         path,
			requestName:  requestName,
			response:     *res,
			responseTime: time,
		})
	}
}
