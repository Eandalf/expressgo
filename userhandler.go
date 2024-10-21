package expressgo

import (
	"io"
	"net/http"
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
	req := &Request{Params: map[string]string{}, Query: map[string]string{}}
	res := &Response{end: false, statusCode: 0, body: ""}
	next := &Next{Next: false, Route: false}
	return req, res, next
}

func (u *UserHandler) setParams(r *http.Request, req *Request) {
	for _, paramsInZone := range u.app.params[r.Pattern] {
		param := ""
		for _, p := range paramsInZone {
			param += p
		}

		values := r.PathValue(param)

		value := ""
		paramIndex := 0
		for _, char := range values {
			if char == '-' {
				if paramIndex+1 < len(paramsInZone) && paramsInZone[paramIndex+1] == "0H" {
					req.Params[paramsInZone[paramIndex]] = value

					// for next param
					value = ""
					paramIndex += 2
				} else {
					value += string(char)
				}
			} else if char == '.' {
				if paramIndex+1 < len(paramsInZone) && paramsInZone[paramIndex+1] == "0D" {
					req.Params[paramsInZone[paramIndex]] = value

					// for next param
					value = ""
					paramIndex += 2
				} else {
					value += string(char)
				}
			} else {
				value += string(char)
			}
		}

		if value != "" && paramIndex < len(paramsInZone) {
			req.Params[paramsInZone[paramIndex]] = value
			paramIndex += 2
		}

		// if any remaining param is not assigned with a value, assign "" to it
		for ; paramIndex < len(paramsInZone); paramIndex += 2 {
			req.Params[paramsInZone[paramIndex]] = ""
		}
	}
}

// Set req.Query[string]string from r.URL.Query().Get(string).
//
// We only accept the first value of a key from the query string.
func (u *UserHandler) setQuery(r *http.Request, req *Request) {
	q := r.URL.Query()
	for k := range q {
		// we only accept the first value of a key from the query string
		req.Query[k] = q.Get(k)
	}
}

// Go through callbacks
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
			if currentCallbackSetIndex+1 > len(u.app.callbacks[u.route])-1 {
				break
			}

			// get the next list of callbacks associated with the designated route
			nextCallbacks := u.app.callbacks[u.route][currentCallbackSetIndex+1]

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
	u.route = r.Pattern

	// prepare custom objects, including req, res, and next
	req, res, next := u.createContext()

	// append params
	u.setParams(r, req)

	// set the query
	u.setQuery(r, req)

	// execute the callbacks
	u.runCallbacks(u.callbacks, 0, req, res, next, w)
}
