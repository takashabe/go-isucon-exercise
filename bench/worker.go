package main

import "net/http"

// Worker is send requests
type Worker struct {
	task   Task
	ctx    Ctx
	result *Result
}

// Need subclass
type Task interface {
	Task()
}

func (w *Worker) getAndCheck(sess *Session, path, requestName string, check func(c *Checker)) {
	req, err := http.NewRequest("GET", path, nil)
	if err != nil {
		return
	}

	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	if sess != nil {
		client.Jar = sess.cookie
	}
	w.requestAndCheck(req, client, requestName, check)
}

func (w *Worker) requestAndCheck(req *http.Request, client *http.Client, requestName string, check func(c *Checker)) {
	PrintDebugf("SEND REQUEST: [%s] %s", requestName, req.URL.Path)
	res, err := client.Do(req)
	if err != nil {
		PrintDebugf("failed to send request %v", err)
		// error is regarded as a server error
		w.result.addResponse(500)
		return
	}

	w.result.addResponse(res.StatusCode)
	if check != nil {
		check(newChecker(w.ctx, w.result))
	}
}
