package main

import (
	"encoding/json"
	"fmt"
	"sync"
)

// Result is save benchmark results
type Result struct {
	Valid        bool             `json:"valid"`
	RequestCount int              `json:"request_count"`
	ElapsedTime  int              `json:"elapsed_time"`
	Response     *ResponseCounter `json:"response"`
	Violations   []*Violation     `json:"violations"`
	mu           sync.Mutex       `json:"-"`
}

func newResult() *Result {
	return &Result{
		Valid:      true,
		Response:   newResponse(),
		Violations: make([]*Violation, 0),
	}
}

func (r *Result) Merge(dst Result) *Result {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.Valid = r.Valid && dst.Valid
	r.RequestCount += dst.RequestCount
	r.ElapsedTime += dst.ElapsedTime

	r.Response.Success += dst.Response.Success
	r.Response.Redirect += dst.Response.Redirect
	r.Response.ClientError += dst.Response.ClientError
	r.Response.ServerError += dst.Response.ServerError
	r.Response.Exception += dst.Response.Exception

	for _, dv := range dst.Violations {
		if rv, ok := r.getViolation(dv.RequestName, dv.Cause); ok {
			rv.Count += dv.Count
			continue
		}
		r.Violations = append(r.Violations, dv)
	}

	return r
}

func (r *Result) addResponse(code int) *Result {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.RequestCount++
	if 200 <= code && code < 300 {
		r.Response.addSuccess()
	} else if 300 <= code && code < 400 {
		r.Response.addRedirect()
	} else if 400 <= code && code < 500 {
		r.Response.addClientError()
	} else {
		r.Response.addServerError()
	}
	return r
}

func (r *Result) addResponseException() *Result {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.RequestCount++
	r.Response.addException()
	return r
}

func (r *Result) addViolation(name, cause string) *Result {
	r.mu.Lock()
	defer r.mu.Unlock()

	if v, ok := r.getViolation(name, cause); ok {
		v.Count++
		return r
	}

	r.Violations = append(r.Violations, &Violation{
		RequestName: name,
		Cause:       cause,
		Count:       1,
	})
	return r
}

func (r *Result) getViolation(name, cause string) (*Violation, bool) {
	for _, v := range r.Violations {
		if v.RequestName == name && v.Cause == cause {
			return v, true
		}
	}
	return nil, false
}

func (r *Result) json() ([]byte, error) {
	return json.MarshalIndent(r, "", "\t")
}

// ResponseCounter holds results for each benchmark
type ResponseCounter struct {
	Success     int        `json:"success"`      // 2xx
	Redirect    int        `json:"redirect"`     // 3xx
	ClientError int        `json:"client_error"` // 4xx
	ServerError int        `json:"server_error"` // 5xx
	Exception   int        `json:"exception"`    // failed request(for example, timeout)
	mu          sync.Mutex `json:"-"`
}

func newResponse() *ResponseCounter {
	return &ResponseCounter{}
}

func (r *ResponseCounter) addSuccess() *ResponseCounter {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.Success++
	return r
}

func (r *ResponseCounter) addRedirect() *ResponseCounter {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.Redirect++
	return r
}

func (r *ResponseCounter) addClientError() *ResponseCounter {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.ClientError++
	return r
}

func (r *ResponseCounter) addServerError() *ResponseCounter {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.ServerError++
	return r
}

func (r *ResponseCounter) addException() *ResponseCounter {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.Exception++
	return r
}

// Violation is save failed requests with cause
type Violation struct {
	RequestName string `json:"request_type"`
	Cause       string `json:"description"`
	Count       int    `json:"num"`
}

func (v *Violation) String() string {
	return fmt.Sprintf("RequestName: %s, Cause: %s, Count: %d", v.RequestName, v.Cause, v.Count)
}
