package main

import (
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
)

// violation text
var (
	causeInvalidResponse      = "パス %s のレスポンスが正しくありません"
	causeStatusCode           = "パス '%s' へのレスポンスコード %d が期待されていましたが %d でした"
	causeRedirectStatusCode   = "レスポンスコードが一時リダイレクトのもの(302, 303, 307)ではなく %d でした"
	causeNoLocation           = "Locationヘッダがありません"
	causeInvalidLocationPath  = "リダイレクト先が %s でなければなりませんが %s でした"
	causeInvalidContentLength = "パス %s に対するレスポンスのサイズが正しくありません: %d bytes"
	causeNoContentLength      = "リクエストパス %s に対して Content-Length がありませんでした"
	causeNoLongerResponse     = "アプリケーションが %d ミリ秒以内に応答しませんでした"
	causeNoStyleSheet         = "スタイルシートのパス %s への参照がありません"
	causeNoNode               = "指定のDOM要素 '%s' が見付かりません"
	causeDifferentNodeCount   = "指定のDOM要素 '%s' が %d 回表示されるはずですが、正しくありません"
	causeFoundNode            = "DOM要素 '%s' は存在しないはずですが、表示されています"
	causeNoContent            = "DOM要素 '%s' で文字列 '%s' を持つものが見付かりません"
	causeDifferentContent     = "DOM要素 '%s' に文字列 '%s' がセットされているはずですが、'%s' となっています"
	causeFoundContent         = "DOM要素 '%s' に文字列 '%s' をもつものは表示されないはずですが、表示されています"
	causeNoBigContent         = "入力されたはずのテキストがDOM要素 '%s' に表示されていません"
	causeNoMatchContent       = "DOM要素 '%s' の中に、テキストが正規表現 '%s' にマッチするものが見つかりません"
	causeDifferentAttribute   = "DOM要素 '%s' のattribute %s の内容が %s になっていません"
)

// document saves goquery.Document for each request type
type document struct {
	doc map[string]*goquery.Document
	mu  sync.Mutex
}

func (d *document) get(key string) (*goquery.Document, bool) {
	d.mu.Lock()
	defer d.mu.Unlock()
	v, ok := d.doc[key]
	return v, ok
}

func (d *document) set(key string, doc *goquery.Document) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.doc[key] = doc
}

// Checker is check benchmark response
type Checker struct {
	ctx          Ctx
	result       *Result
	path         string
	requestName  string
	response     http.Response
	responseTime time.Duration
	document     document
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
		c.addViolation(causeNoLocation)
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

func (c *Checker) isContentLength(size int) {
	cl := c.response.Header.Get("Content-Length")
	if cl == "" {
		c.addViolation(causeNoContentLength)
		return
	}
	if i, err := strconv.Atoi(cl); err != nil || i != size {
		c.addViolation(fmt.Sprintf(causeInvalidContentLength, c.path, i))
		return
	}
}

func (c *Checker) respondUntil(limit time.Duration) {
	if c.responseTime > limit {
		return
	}
	c.addViolation(fmt.Sprintf(causeNoLongerResponse, limit))
}

func (c *Checker) getDocument() (*goquery.Document, error) {
	if c.document.doc == nil {
		c.document.doc = make(map[string]*goquery.Document)
	} else {
		if d, ok := c.document.get(c.requestName); ok {
			return d, nil
		}
	}

	d, err := goquery.NewDocumentFromResponse(&c.response)
	if err != nil {
		c.addViolation(fmt.Sprintf(causeInvalidResponse, c.path))
		return nil, err
	}
	c.document.set(c.requestName, d)
	return d, nil
}

func (c *Checker) hasStyleSheet(path string) {
	doc, err := c.getDocument()
	if err != nil {
		return
	}

	link := doc.Find("link")
	if rel, ok := link.Attr("rel"); ok && rel == "stylesheet" {
		if href, ok := link.Attr("href"); ok && href == path {
			// OK
			return
		}
	}
	c.addViolation(fmt.Sprintf(causeNoStyleSheet, path))
}

func (c *Checker) hasNode(selector string) {
	doc, err := c.getDocument()
	if err != nil {
		return
	}

	if doc.Find(selector).Size() > 0 {
		// OK
		return
	}
	c.addViolation(fmt.Sprintf(causeNoNode, selector))
}

func (c *Checker) nodeCount(selector string, num int) {
	doc, err := c.getDocument()
	if err != nil {
		return
	}

	if doc.Find(selector).Size() == num {
		// OK
		return
	}
	c.addViolation(fmt.Sprintf(causeDifferentNodeCount, selector, num))
}

func (c *Checker) missingNode(selector string) {
	doc, err := c.getDocument()
	if err != nil {
		return
	}

	if doc.Find(selector).Size() <= 0 {
		// OK
		return
	}
	c.addViolation(fmt.Sprintf(causeFoundNode, selector))
}

func (c *Checker) hasContent(selector, text string) {
	doc, err := c.getDocument()
	if err != nil {
		return
	}

	nodes := []*goquery.Selection{}
	doc.Find(selector).Each(func(i int, s *goquery.Selection) {
		if s.Text() != "" {
			nodes = append(nodes, s)
		}
	})
	if len(nodes) == 0 {
		c.addViolation(fmt.Sprintf(causeNoContent, selector, text))
		return
	}
	for _, s := range nodes {
		if s.Text() == text {
			// OK
			return
		}
	}
	c.addViolation(fmt.Sprintf(causeDifferentContent, selector, text, nodes[0].Text()))
}

func (c *Checker) missingContent(selector, text string) {
	doc, err := c.getDocument()
	if err != nil {
		return
	}

	nodes := []*goquery.Selection{}
	doc.Find(selector).Each(func(i int, s *goquery.Selection) {
		if s.Text() == text {
			nodes = append(nodes, s)
		}
	})
	if len(nodes) == 0 {
		// OK
		return
	}
	c.addViolation(fmt.Sprintf(causeFoundContent, selector, text))
}

var trimBR = regexp.MustCompile(`(?m)<(br|BR|Br|bR) */?>`)

func (c *Checker) hasBigContent(selector, text string) {
	doc, err := c.getDocument()
	if err != nil {
		return
	}

	text = strings.Join(strings.Split(text, "\n"), "")
	nodes := []*goquery.Selection{}
	doc.Find(selector).Each(func(i int, s *goquery.Selection) {
		if trimBR.ReplaceAllString(s.Text(), "") == text {
			nodes = append(nodes, s)
			return
		}
	})
	if len(nodes) > 0 {
		// OK
		return
	}
	c.addViolation(fmt.Sprintf(causeNoBigContent, selector))
}

func (c *Checker) matchContent(selector, regex string) {
	reg, err := regexp.Compile(regex)
	if err != nil {
		c.addViolation(fmt.Sprintf(causeNoMatchContent, selector, regex))
		return
	}
	doc, err := c.getDocument()
	if err != nil {
		return
	}

	nodes := []*goquery.Selection{}
	doc.Find(selector).Each(func(i int, s *goquery.Selection) {
		if reg.MatchString(s.Text()) {
			nodes = append(nodes, s)
			return
		}
	})
	if len(nodes) > 0 {
		// OK
		return
	}
	c.addViolation(fmt.Sprintf(causeNoMatchContent, selector, regex))
}

func (c *Checker) contentFunc(selector, cause string, f func(s *goquery.Selection) bool) {
	doc, err := c.getDocument()
	if err != nil {
		return
	}

	nodes := []*goquery.Selection{}
	doc.Find(selector).Each(func(i int, s *goquery.Selection) {
		if f(s) {
			nodes = append(nodes, s)
			return
		}
	})
	if len(nodes) > 0 {
		// OK
		return
	}
	c.addViolation(cause)
}

func (c *Checker) attribute(selector, attr, text string) {
	doc, err := c.getDocument()
	if err != nil {
		return
	}

	nodes := []*goquery.Selection{}
	doc.Find(selector).Each(func(i int, s *goquery.Selection) {
		if t, ok := s.Attr(attr); ok && t == text {
			nodes = append(nodes, s)
			return
		}
	})
	if len(nodes) > 0 {
		// OK
		return
	}
	c.addViolation(fmt.Sprintf(causeDifferentAttribute, selector, attr, text))
}
