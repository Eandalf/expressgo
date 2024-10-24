package expressgo

import "net/http"

type Request struct {
	native *http.Request
	Params map[string]string
	Query  map[string]string
}
