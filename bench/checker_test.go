package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/PuerkitoBio/goquery"
)

func testChecker() *Checker {
	return &Checker{
		ctx:         defaultCtx,
		result:      newResult(),
		path:        "/",
		requestName: "TEST",
		response:    *httptest.NewRecorder().Result(),
	}
}

func testResponse(code int) *http.Response {
	recorder := httptest.NewRecorder()
	recorder.WriteHeader(code)
	return recorder.Result()
}

func testRedirectResponse(path string, code int) *http.Response {
	recorder := httptest.NewRecorder()
	recorder.Header().Set("Location", path)
	recorder.WriteHeader(code)
	return recorder.Result()
}

func testContentLengthResponse(length int) *http.Response {
	recorder := httptest.NewRecorder()
	recorder.Header().Set("Content-Length", fmt.Sprint(length))
	return recorder.Result()
}

func testBodyResponse(body []byte) *http.Response {
	recorder := httptest.NewRecorder()
	recorder.Body.Write(body)
	return recorder.Result()
}

func TestIsStatusCode(t *testing.T) {
	cases := []struct {
		input        int
		expectResult *Result
	}{
		{
			200,
			newResult(),
		},
		{
			500,
			newResult().addViolation("TEST", fmt.Sprintf(causeStatusCode, "/", 500, 200)),
		},
	}
	for i, c := range cases {
		checker := testChecker()
		checker.isStatusCode(c.input)
		if !reflect.DeepEqual(checker.result, c.expectResult) {
			t.Errorf("#%d: want %v, got %v", i, c.expectResult, checker.result)
		}
	}
}

func TestIsRedirect(t *testing.T) {
	cases := []struct {
		response     *http.Response
		input        string
		expectResult *Result
	}{
		{
			testRedirectResponse("/test", 302),
			"/test",
			newResult(),
		},
		{
			testRedirectResponse("/test", 200),
			"/test",
			newResult().addViolation("TEST", fmt.Sprintf(causeRedirectStatusCode, 200)),
		},
		{
			testResponse(302),
			"/test",
			newResult().addViolation("TEST", fmt.Sprintf(causeNoLocation)),
		},
		{
			testRedirectResponse("http://localhost/test", 302),
			"/test",
			newResult(),
		},
		{
			testRedirectResponse("http://localhost/foo", 302),
			"/test",
			newResult().addViolation("TEST", fmt.Sprintf(causeInvalidLocationPath, "/test", "/foo")),
		},
	}
	for i, c := range cases {
		checker := testChecker()
		checker.response = *c.response
		checker.isRedirect(c.input)
		if !reflect.DeepEqual(checker.result, c.expectResult) {
			t.Errorf("#%d: want %v, got %v", i, c.expectResult, checker.result)
		}
	}
}

func TestIsContentLength(t *testing.T) {
	cases := []struct {
		response     *http.Response
		input        int
		expectResult *Result
	}{
		{
			testContentLengthResponse(1),
			1,
			newResult(),
		},
		{
			testContentLengthResponse(-1),
			1,
			newResult().addViolation("TEST", fmt.Sprintf(causeInvalidContentLength, "/", -1)),
		},
		{
			testResponse(200),
			1,
			newResult().addViolation("TEST", causeNoContentLength),
		},
	}
	for i, c := range cases {
		checker := testChecker()
		checker.response = *c.response
		checker.isContentLength(c.input)
		if !reflect.DeepEqual(checker.result, c.expectResult) {
			t.Errorf("#%d: want %v, got %v", i, c.expectResult, checker.result)
		}
	}
}

func TestHasStyleSheet(t *testing.T) {
	cases := []struct {
		response     *http.Response
		input        string
		expectResult *Result
	}{
		{
			testBodyResponse([]byte(`
<head>
  <link rel="stylesheet" href="foo">
</head>`)),
			"foo",
			newResult(),
		},
		{
			testBodyResponse([]byte(`
<head>
  <link rel="stylesheet" href="foo">
</head>`)),
			"bar",
			newResult().addViolation("TEST", fmt.Sprintf(causeNoStyleSheet, "bar")),
		},
	}
	for i, c := range cases {
		checker := testChecker()
		checker.response = *c.response
		// goquery require response.request
		checker.response.Request = httptest.NewRequest("GET", "/", nil)
		checker.hasStyleSheet(c.input)
		if !reflect.DeepEqual(checker.result, c.expectResult) {
			t.Errorf("#%d: want %v, got %v", i, c.expectResult, checker.result)
		}
	}
}

func TestHasNode(t *testing.T) {
	cases := []struct {
		input        string
		expectResult *Result
	}{
		{
			".bar",
			newResult(),
		},
		{
			"dd.foo div.bar #foobar",
			newResult(),
		},
		{
			"#bar",
			newResult().addViolation("TEST", fmt.Sprintf(causeNoNode, "#bar")),
		},
	}
	for i, c := range cases {
		response := testBodyResponse([]byte(`
<dd>
  <dd class='foo'>
    <div class='bar'>
      <div id='foobar'>
    </div>
  </dd>
</dd>`))
		checker := testChecker()
		checker.response = *response
		// goquery require response.request
		checker.response.Request = httptest.NewRequest("GET", "/", nil)
		checker.hasNode(c.input)
		if !reflect.DeepEqual(checker.result, c.expectResult) {
			t.Errorf("#%d: want %v, got %v", i, c.expectResult, checker.result)
		}
	}
}

func TestNodeCount(t *testing.T) {
	cases := []struct {
		inputSelector string
		inputNum      int
		expectResult  *Result
	}{
		{
			".foo",
			2,
			newResult(),
		},
		{
			"div",
			3,
			newResult(),
		},
		{
			".foobar",
			1,
			newResult().addViolation("TEST", fmt.Sprintf(causeDifferentNodeCount, ".foobar", 1)),
		},
	}
	for i, c := range cases {
		response := testBodyResponse([]byte(`
<div class='foo' />
<div class='foo' />
<div class='bar' />`))
		checker := testChecker()
		checker.response = *response
		// goquery require response.request
		checker.response.Request = httptest.NewRequest("GET", "/", nil)
		checker.nodeCount(c.inputSelector, c.inputNum)
		if !reflect.DeepEqual(checker.result, c.expectResult) {
			t.Errorf("#%d: want %v, got %v", i, c.expectResult, checker.result)
		}
	}
}

func TestMissingNode(t *testing.T) {
	cases := []struct {
		input        string
		expectResult *Result
	}{
		{
			".foobar",
			newResult(),
		},
		{
			"div",
			newResult().addViolation("TEST", fmt.Sprintf(causeFoundNode, "div")),
		},
	}
	for i, c := range cases {
		response := testBodyResponse([]byte(`
<div class='foo' />
<div class='foo' />
<div class='bar' />`))
		checker := testChecker()
		checker.response = *response
		// goquery require response.request
		checker.response.Request = httptest.NewRequest("GET", "/", nil)
		checker.missingNode(c.input)
		if !reflect.DeepEqual(checker.result, c.expectResult) {
			t.Errorf("#%d: want %v, got %v", i, c.expectResult, checker.result)
		}
	}
}

func TestHasContent(t *testing.T) {
	cases := []struct {
		inputSelector string
		inputText     string
		expectResult  *Result
	}{
		{
			".foo",
			"foo",
			newResult(),
		},
		{
			"#foo",
			"foo",
			newResult().addViolation("TEST", fmt.Sprintf(causeNoContent, "#foo", "foo")),
		},
		{
			".bar",
			"bar",
			newResult().addViolation("TEST", fmt.Sprintf(causeDifferentContent, ".bar", "bar", "foo")),
		},
	}
	for i, c := range cases {
		response := testBodyResponse([]byte(`
<div class='foo' />
<div class='foo'>foo</div>
<div class='bar'>foo</div>`))
		checker := testChecker()
		checker.response = *response
		// goquery require response.request
		checker.response.Request = httptest.NewRequest("GET", "/", nil)
		checker.hasContent(c.inputSelector, c.inputText)
		if !reflect.DeepEqual(checker.result, c.expectResult) {
			t.Errorf("#%d: want %v, got %v", i, c.expectResult, checker.result)
		}
	}
}

func TestMissingContent(t *testing.T) {
	cases := []struct {
		inputSelector string
		inputText     string
		expectResult  *Result
	}{
		{
			".bar",
			"bar",
			newResult(),
		},
		{
			".foo",
			"foo",
			newResult().addViolation("TEST", fmt.Sprintf(causeFoundContent, ".foo", "foo")),
		},
	}
	for i, c := range cases {
		response := testBodyResponse([]byte(`
<div class='foo' />
<div class='foo'>foo</div>
<div class='bar'>foo</div>`))
		checker := testChecker()
		checker.response = *response
		// goquery require response.request
		checker.response.Request = httptest.NewRequest("GET", "/", nil)
		checker.missingContent(c.inputSelector, c.inputText)
		if !reflect.DeepEqual(checker.result, c.expectResult) {
			t.Errorf("#%d: want %v, got %v", i, c.expectResult, checker.result)
		}
	}
}

func TestHasBigContent(t *testing.T) {
	cases := []struct {
		inputSelector string
		inputText     string
		expectResult  *Result
	}{
		{
			".foo",
			"foo\nbar\nfoo\n\nbar\n",
			newResult(),
		},
		{
			".foo",
			"foobarfoobar",
			newResult(),
		},
		{
			".foo",
			"foobarfoo",
			newResult().addViolation("TEST", fmt.Sprintf(causeNoBigContent, ".foo")),
		},
	}
	for i, c := range cases {
		response := testBodyResponse([]byte(`
<div class='foo'>foo<br >bar<BR/>foo<Br    ><bR       />bar<br></div>`))
		checker := testChecker()
		checker.response = *response
		// goquery require response.request
		checker.response.Request = httptest.NewRequest("GET", "/", nil)
		checker.hasBigContent(c.inputSelector, c.inputText)
		if !reflect.DeepEqual(checker.result, c.expectResult) {
			t.Errorf("#%d: want %v, got %v", i, c.expectResult, checker.result)
		}
	}
}

func TestMatchContent(t *testing.T) {
	cases := []struct {
		inputSelector string
		inputRegex    string
		expectResult  *Result
	}{
		{
			".foo",
			"foo.*bar$",
			newResult(),
		},
		{
			".foo",
			"none",
			newResult().addViolation("TEST", fmt.Sprintf(causeNoMatchContent, ".foo", "none")),
		},
		{
			".foo",
			"(",
			newResult().addViolation("TEST", fmt.Sprintf(causeNoMatchContent, ".foo", "(")),
		},
	}
	for i, c := range cases {
		response := testBodyResponse([]byte(`
<div class='foo'>foo,123456789,bar</div>`))
		checker := testChecker()
		checker.response = *response
		// goquery require response.request
		checker.response.Request = httptest.NewRequest("GET", "/", nil)
		checker.matchContent(c.inputSelector, c.inputRegex)
		if !reflect.DeepEqual(checker.result, c.expectResult) {
			t.Errorf("#%d: want %v, got %v", i, c.expectResult, checker.result)
		}
	}
}

func TestContentFunc(t *testing.T) {
	cases := []struct {
		inputSelector string
		inputCause    string
		inputFunc     func(s *goquery.Selection) bool
		expectResult  *Result
	}{
		{
			".foo",
			"test",
			func(s *goquery.Selection) bool {
				return s.Text() == "foo"
			},
			newResult(),
		},
		{
			".foo",
			"test",
			func(s *goquery.Selection) bool {
				return false
			},
			newResult().addViolation("TEST", "test"),
		},
	}
	for i, c := range cases {
		response := testBodyResponse([]byte(`
<div class='foo'>foo</div>`))
		checker := testChecker()
		checker.response = *response
		// goquery require response.request
		checker.response.Request = httptest.NewRequest("GET", "/", nil)
		checker.contentFunc(c.inputSelector, c.inputCause, c.inputFunc)
		if !reflect.DeepEqual(checker.result, c.expectResult) {
			t.Errorf("#%d: want %v, got %v", i, c.expectResult, checker.result)
		}
	}
}

func TestAttribute(t *testing.T) {
	cases := []struct {
		inputSelector string
		inputAttr     string
		inputText     string
		expectResult  *Result
	}{
		{
			"a",
			"href",
			"bar",
			newResult(),
		},
		{
			".foo",
			"attr",
			"div",
			newResult(),
		},
		{
			".foo",
			"href",
			"foo",
			newResult().addViolation("TEST", fmt.Sprintf(causeDifferentAttribute, ".foo", "href", "foo")),
		},
	}
	for i, c := range cases {
		response := testBodyResponse([]byte(`
<div class='foo' attr='div'>
  <a href='foo'>test1</a>
  <a href='bar'>test2</a>
    <div attr='div'>nested</div>
  <a href='bar'>test3</a>
</div>`))
		checker := testChecker()
		checker.response = *response
		// goquery require response.request
		checker.response.Request = httptest.NewRequest("GET", "/", nil)
		checker.attribute(c.inputSelector, c.inputAttr, c.inputText)
		if !reflect.DeepEqual(checker.result, c.expectResult) {
			t.Errorf("#%d: want %v, got %v", i, c.expectResult, checker.result)
		}
	}
}

func TestMultipleCallDocument(t *testing.T) {
	response := testBodyResponse([]byte(`
<div id="login-form">
  <form method="POST" action="/login">
    <div class="col-md-4 input-group">
      <span class="input-group-addon">E-mail</span>
      <input class="form-control" type="text" name="email" placeholder="E-mail address" />
    </div>
    <div class="col-md-4 input-group">
      <span class="input-group-addon">パスワード</span>
      <input class="form-control" type="password" name="password" />
    </div>
    <div class="col-md-1 input-group">
      <input class="btn btn-default" type="submit" name="Login" value="Login" />
    </div>
  </form>
</div>`))
	checker := testChecker()
	checker.response = *response
	checker.response.Request = httptest.NewRequest("GET", "/", nil)

	args := []struct {
		selector string
		num      int
	}{
		{"form input[type=password]", 1},
		{"form input[type=submit]", 1},
		{"form input[type=text]", 1},
	}
	for i, a := range args {
		checker.nodeCount(a.selector, a.num)
		if len(checker.result.Violations) != 0 {
			t.Errorf("#%d: want no violations, got violations %v", i, checker.result.Violations)
		}
	}
}
