package expressgo

import (
	"io"
	"net/http"
	"strings"
)

type Next struct {
	Next  bool
	Route bool
}

type Callback func(req *Request, res *Response, next *Next)

type UserHandler struct {
	app       *App
	callbacks []Callback
	route     string
}

func (u *UserHandler) createContext() (*Request, *Response, *Next) {
	req := &Request{Params: make(map[string]string), Query: make(map[string]string)}
	res := &Response{end: false, statusCode: 0, body: ""}
	next := &Next{Next: false, Route: false}
	return req, res, next
}

// go through callbacks
func (u *UserHandler) runCallbacks(
	callbacks []Callback,
	currentCallbackSetIndex int,
	req *Request,
	res *Response,
	next *Next,
	w http.ResponseWriter,
) {
	for _, c := range callbacks {
		c(req, res, next)

		// perform the write, res -> ResponseWriter
		if res.statusCode != 0 {
			w.WriteHeader(res.statusCode)
		}
		if res.body != "" {
			io.WriteString(w, res.body)
		}

		// check next route
		if next.Route {
			// ensure the index is not out of the boundary
			if currentCallbackSetIndex+1 > len(u.app.routes[u.route])-1 {
				break
			}

			// get the next set of callbacks associated with the designated path
			nextCallbacks := u.app.routes[u.route][currentCallbackSetIndex+1]

			// run callbacks
			u.runCallbacks(
				nextCallbacks,
				currentCallbackSetIndex+1,
				req,
				res,
				&Next{Next: false, Route: false},
				w,
			)
			break
		}

		// check next status
		if !next.Next || res.end {
			break
		}
	}
}

func (u *UserHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// save the route first
	u.route = strings.TrimSuffix(r.Method+" "+r.URL.Path, "/")

	// prepare custom objects, including req, res, and next
	req, res, next := u.createContext()

	// execute the callbacks
	u.runCallbacks(u.callbacks, 0, req, res, next, w)
}
