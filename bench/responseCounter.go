package main

// ResponseCounter holds results for each benchmark
type ResponseCounter struct {
	success     int // 2xx
	redirect    int // 3xx
	clientError int // 4xx
	serverError int // 5xx
	exception   int // failed request(for example, timeout)
}

func newResponse() *ResponseCounter {
	return &ResponseCounter{}
}

func (r *ResponseCounter) addSuccess() *ResponseCounter {
	r.success++
	return r
}

func (r *ResponseCounter) addRedirect() *ResponseCounter {
	r.redirect++
	return r
}

func (r *ResponseCounter) addClientError() *ResponseCounter {
	r.clientError++
	return r
}

func (r *ResponseCounter) addServerError() *ResponseCounter {
	r.serverError++
	return r
}

func (r *ResponseCounter) addException() *ResponseCounter {
	r.exception++
	return r
}
