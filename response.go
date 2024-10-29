package expressgo

import (
	"net/http"
	"strings"
)

type Response struct {
	native     http.ResponseWriter
	end        bool
	statusCode int
	body       string
}

// Stop further writes to the response.
func (res *Response) End() {
	res.end = true
}

// Add a value to a response header, field: value. The field is case-insensitive.
func (res *Response) Append(field string, value string) {
	// if end is already designated, this method should be a no-op
	if res.end {
		return
	}

	res.native.Header().Add(field, value)
}

// Set a response header, field: value. The field is case-insensitive.
func (res *Response) Set(field string, value string) {
	// if end is already designated, this method should be a no-op
	if res.end {
		return
	}

	res.native.Header().Set(field, value)
}

// Get a response header specified by the field. The field is case-insensitive.
func (res *Response) Get(field string) string {
	// This implementation is based on Mozilla's mozilla-central/netwerk/protocol/http/nsHttpHeaderArray.h
	// https://github.com/bnoordhuis/mozilla-central/blob/master/netwerk/protocol/http/nsHttpHeaderArray.h#L185
	specialHeader := [...]string{
		"Set-Cookie",
		"WWW-Authenticate",
		"Proxy-Authenticate",
	}

	for _, h := range specialHeader {
		// case-insensite comparison
		if strings.EqualFold(field, h) {
			return res.native.Header().Get(field)
		}
	}

	values := res.native.Header().Values(field)
	return strings.Join(values, ",")
}

// Send the response.
func (res *Response) Send(body string) {
	// if end is already designated, this method should be a no-op
	if res.end {
		return
	}

	res.body = body
	res.end = true
}

// Send the response with a status code.
func (res *Response) SendStatus(statusCode int) {
	// if end is already designated, this method should be a no-op
	if res.end {
		return
	}

	res.statusCode = statusCode
	res.end = true
}

// Set the HTTP status code of the response, it is chainable.
func (res *Response) Status(statusCode int) *Response {
	res.statusCode = statusCode
	return res
}
