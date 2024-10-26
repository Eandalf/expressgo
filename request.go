package expressgo

import "net/http"

type Request struct {
	Native *http.Request
	Params map[string]string
	Query  map[string]string
	Body   interface{}
	err    error
}

type BodyJsonBase map[string]interface{}
