package expressgo

import (
	"encoding/json"
	"net/http"
	"strings"
)

type Request struct {
	Native *http.Request
	Params map[string]string
	Query  map[string]string
	Body   interface{}
	err    error
}

type BodyJsonBase map[string]json.RawMessage

type BodyFormUrlEncoded map[string]string

// Get a request header specified by the field. The field is case-insensitive.
func (req *Request) Get(field string) string {
	values := req.Native.Header.Values(field)
	return strings.Join(values, ",")
}

// Alias of req.Get(string).
func (req *Request) Header(field string) string {
	return req.Get(field)
}
