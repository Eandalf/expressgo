package expressgo

import "net/http"

type Request struct {
	Native *http.Request
	Params map[string]string
	Query  map[string]string
}
