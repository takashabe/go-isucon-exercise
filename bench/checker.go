package main

import (
	"fmt"
	"net/http"
)

type Checker struct {
	ctx         Ctx
	result      *Result
	requestName string
	response    http.Response
}

func (c *Checker) isStatusCode(code int) {
	if c.response.StatusCode != code {
		c.result.addViolation(c.requestName, fmt.Sprintf("パス '%s' へのレスポンスコード %d が期待されていましたが %d でした", c.response.Request.URL.Path, code, c.response.StatusCode))
	}
}
