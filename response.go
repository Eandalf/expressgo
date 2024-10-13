package expressgo

type Response struct {
	end        bool
	statusCode int
	body       string
}

func (res *Response) End() {
	res.end = true
}

func (res *Response) Send(body string) {
	// if end is already designated, this method should be a no-op
	if res.end {
		return
	}

	res.body = body
	res.end = true
}

func (res *Response) SendStatus(statusCode int) {
	// if end is already designated, this method should be a no-op
	if res.end {
		return
	}

	res.statusCode = statusCode
	res.end = true
}

// chainable
func (res *Response) Status(statusCode int) *Response {
	res.statusCode = statusCode
	return res
}
