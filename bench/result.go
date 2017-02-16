package main

// Result is save benchmark results
type Result struct {
	valid        bool
	requestCount int
	elapsedTime  int
	response     *ResponseCounter
	violations   []*Violation
}

// Violation is save failed requests with cause
type Violation struct {
	requestName string
	cause       string
	count       int
}

func newResult() *Result {
	return &Result{
		response:   newResponse(),
		violations: make([]*Violation, 0),
	}
}

func (r *Result) Merge(dst Result) *Result {
	r.valid = r.valid && dst.valid
	r.requestCount += dst.requestCount
	r.elapsedTime += dst.elapsedTime

	r.response.success += dst.response.success
	r.response.redirect += dst.response.redirect
	r.response.clientError += dst.response.clientError
	r.response.serverError += dst.response.serverError
	r.response.exception += dst.response.exception

	for _, dv := range dst.violations {
		if rv, ok := r.getViolation(dv.requestName, dv.cause); ok {
			rv.count += dv.count
			continue
		}
		r.violations = append(r.violations, dv)
	}

	return r
}

func (r *Result) addResponse(code int) *Result {
	r.requestCount++
	if 200 <= code && code < 300 {
		r.response.addSuccess()
	} else if 300 <= code && code < 400 {
		r.response.addRedirect()
	} else if 400 <= code && code < 500 {
		r.response.addClientError()
	} else {
		r.response.addServerError()
	}
	return r
}

func (r *Result) addResponseException() *Result {
	r.requestCount++
	r.response.addException()
	return r
}

func (r *Result) addViolation(name, cause string) *Result {
	if v, ok := r.getViolation(name, cause); ok {
		v.count++
		return r
	}

	r.violations = append(r.violations, &Violation{
		requestName: name,
		cause:       cause,
		count:       1,
	})
	return r
}

func (r *Result) getViolation(name, cause string) (*Violation, bool) {
	for _, v := range r.violations {
		if v.requestName == name && v.cause == cause {
			return v, true
		}
	}
	return nil, false
}
