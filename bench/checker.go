package main

import (
	"fmt"
	"net/http"
)

// violation text
var (
	causeStatusCode          = "パス '%s' へのレスポンスコード %d が期待されていましたが %d でした"
	causeRedirectStatusCode  = "レスポンスコードが一時リダイレクトのもの(302, 303, 307)ではなく %d でした"
	causeNoneLocation        = "Locationヘッダがありません"
	causeInvalidLocationPath = "リダイレクト先が %s でなければなりませんが %s でした"
)

type Checker struct {
	ctx         Ctx
	result      *Result
	path        string
	requestName string
	response    http.Response
}

func (c *Checker) statusCode() int {
	return c.response.StatusCode
}

func (c *Checker) addViolation(cause string) {
	c.result.addViolation(c.requestName, cause)
}

func (c *Checker) hasViolation() bool {
	return len(c.result.Violations) > 0
}

func (c *Checker) isStatusCode(code int) {
	if c.statusCode() != code {
		c.addViolation(fmt.Sprintf(causeStatusCode, c.path, code, c.statusCode()))
	}
}

func (c *Checker) isRedirect(path string) {
	// check HTTP status code
	wantStatusCode := []int{302, 303, 307}
	isValidStatusCode := false
	for _, v := range wantStatusCode {
		if v == c.statusCode() {
			isValidStatusCode = true
			break
		}
	}
	if !isValidStatusCode {
		c.addViolation(fmt.Sprintf(causeRedirectStatusCode, c.statusCode()))
		return
	}

	// check location header
	loc, err := c.response.Location()
	if err != nil {
		c.addViolation(causeNoneLocation)
		return
	}

	// check url
	if loc.String() == c.ctx.uri(path) {
		// OK
		return
	}
	// check url(other than port)
	if loc.Host == "" || loc.Host == c.ctx.host && loc.Path == path {
		// OK
		return
	}
	c.addViolation(fmt.Sprintf(causeInvalidLocationPath, path, loc.Path))
}
