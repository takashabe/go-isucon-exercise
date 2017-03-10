package main

import (
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/pkg/errors"
)

type Driver struct {
	result *Result
	ctx    Ctx
}

func (d *Driver) get(sess *Session, path string) {
	d.getAndCheck(sess, path, "", nil)
}

func (d *Driver) getAndCheck(sess *Session, path, requestName string, check func(c *Checker)) {
	req, err := http.NewRequest("GET", d.ctx.uri(path), nil)
	if err != nil {
		PrintDebugf("failed to generate request %v", err)
		// error is regarded as a client error
		d.result.addResponse(400)
		return
	}

	// TODO: reuse global defined http.Client (must reuse transport)
	client := &http.Client{
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

	// TODO: reuse global defined http.Client (must reuse transport)
	client := &http.Client{
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
	res, err := client.Do(req)
	if err != nil {
		PrintDebugf("failed to send request. path=%s, error=%v", path, err)
		// error is regarded as a server error
		d.result.addResponse(500)
		return
	}

	d.result.addResponse(res.StatusCode)
	if check != nil {
		check(&Checker{
			ctx:         d.ctx,
			result:      d.result,
			path:        path,
			requestName: requestName,
			response:    *res,
		})
	}
}
