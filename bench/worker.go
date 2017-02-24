package main

import "net/http"

// Worker is send requests
// TODO: export check request functions from Worker
type Worker struct {
	ctx    Ctx
	result *Result
}

func newWorker() *Worker {
	return &Worker{
		ctx:    *newCtx(),
		result: newResult(),
	}
}

func (w *Worker) getAndCheck(sess *Session, path, requestName string, check func(c *Checker)) {
	req, err := http.NewRequest("GET", w.ctx.uri(path), nil)
	if err != nil {
		// TODO: error handling
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
	w.requestAndCheck(path, requestName, req, client, check)
}

func (w *Worker) requestAndCheck(path, requestName string, req *http.Request, client *http.Client, check func(c *Checker)) {
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
		check(&Checker{
			ctx:         w.ctx,
			result:      w.result,
			path:        path,
			requestName: requestName,
			response:    *res,
		})
	}
}
