package expressgo

import (
	"fmt"
	"io"
	"net/http"
)

type Next struct {
	Next  bool
	Route bool
	Err   error
}

type Callback func(req *Request, res *Response, next *Next)

type ErrorCallback func(err error, req *Request, res *Response, next *Next)

type UserHandler struct {
	app       *App
	callbacks []Callback
	route     string
}

func (u *UserHandler) createContext(r *http.Request, w http.ResponseWriter) (*Request, *Response) {
	req := &Request{
		Native: r,
		Params: map[string]string{},
		Query:  map[string]string{},
	}
	res := &Response{
		native:     w,
		end:        false,
		statusCode: 0,
		body:       "",
	}
	return req, res
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

// Run a callback, for recovering an error
func (u *UserHandler) runCallback(
	c Callback,
	req *Request,
	res *Response,
	next *Next,
) {
	// recover from panic of callbacks
	defer func() {
		// do not recover in development environments
		if appEnv, ok := u.app.GetData(configKeyAppEnv).(string); ok && appEnv != "development" {
			if r := recover(); r != nil {
				next.Err = fmt.Errorf("%#v", r)
			}
		}
	}()

	c(req, res, next)
}

// Go through callbacks
func (u *UserHandler) runCallbacks(
	callbacks []Callback,
	currentCallbackSetIndex int,
	req *Request,
	res *Response,
	w http.ResponseWriter,
) {
	for pos, c := range callbacks {
		// create a new next for each callback
		next := &Next{Next: false, Route: false, Err: nil}

		u.runCallback(c, req, res, next)

		// perform the write, res -> ResponseWriter
		if res.statusCode != 0 {
			w.WriteHeader(res.statusCode)
		}
		if res.body != "" {
			io.WriteString(w, res.body)
		}

		// transfer the error from next to req
		if next.Err != nil {
			req.err = next.Err
			next.Err = nil
		}
		// if the error is not consumed, activate next.Next or next.Route to pass the error to error handlers down the callback lists
		if req.err != nil {
			if pos == (len(callbacks) - 1) {
				next.Route = true
			} else {
				next.Next = true
			}
		}

		// do not proceed if the respond is meant to be sent, even with next.Next ot next.Route is set
		if res.end {
			break
		}

		// next.Next takes precedence over next.Route
		// next.Next is meaningless with the last callback in the current callback list
		// this check is implemented to have next.Route in effect with the last callback
		if next.Next && pos != (len(callbacks)-1) {
			continue
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
				w,
			)
			break
		}

		// check next status
		if !next.Next {
			break
		}
	}
}

func (u *UserHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// save the route first
	u.route = r.Pattern

	// prepare custom objects, including req, res, and next
	req, res := u.createContext(r, w)

	// append params
	u.setParams(r, req)

	// set the query
	u.setQuery(r, req)

	// execute the callbacks
	u.runCallbacks(u.callbacks, 0, req, res, w)
}
