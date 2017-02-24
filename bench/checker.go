package main

import (
	"fmt"
	"net/http"
	"net/url"
)

type Checker struct {
	ctx         Ctx
	result      *Result
	path        string
	requestName string
	response    http.Response
}

func (c *Checker) getStatusCode() int {
	return c.response.StatusCode
}

func (c *Checker) addViolation(cause string) {
	c.result.addViolation(c.requestName, cause)
}

func (c *Checker) hasViolation() bool {
	return len(c.result.Violations) > 0
}

func (c *Checker) isStatusCode(code int) {
	if c.getStatusCode() != code {
		c.addViolation(fmt.Sprintf("パス '%s' へのレスポンスコード %d が期待されていましたが %d でした", c.path, code, c.getStatusCode()))
	}
}

func (c *Checker) isRedirect(path string) {
	// check HTTP status code
	wantStatusCode := []int{302, 303, 307}
	isValidStatusCode := false
	for _, v := range wantStatusCode {
		if v == c.getStatusCode() {
			isValidStatusCode = true
			break
		}
	}
	if isValidStatusCode {
		c.addViolation(fmt.Sprintf("レスポンスコードが一時リダイレクトのもの(302, 303, 307)ではなく %d でした", c.getStatusCode()))
		return
	}

	// check location header
	loc, err := c.response.Location()
	if err != nil {
		c.addViolation("Locationヘッダがありません")
	} else if loc.Path == c.ctx.uri(path) {
		// pass the check
		return
	}

	// check url format
	if url, err := url.Parse(path); err == nil {
		if url.Host == "" || url.Host == c.ctx.host && url.Path == path {
			// pass the check
			return
		}
	}
	c.addViolation(fmt.Sprintf("リダイレクト先が %s でなければなりませんが %s でした", path, loc.Path))
}
