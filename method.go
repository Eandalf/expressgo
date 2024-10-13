package expressgo

import (
	"io"
	"net/http"
)

type Next struct {
	next  bool
	route string
}

type Callback func(req *Request, res *Response, next *Next)

type UserHandler struct {
	app       *App
	callbacks []Callback
}

func (u *UserHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// prepare custom objects, including req, res, and next
	req, res, next := u.app.createConnection()

	// go through callbacks
	for _, c := range u.callbacks {
		c(req, res, next)

		// perform the write, res -> ResponseWriter
		if res.statusCode != 0 {
			w.WriteHeader(res.statusCode)
		}
		if res.body != "" {
			io.WriteString(w, res.body)
		}

		// check next status
		if !next.next || res.end {
			break
		} else if next.route != "" {
			// not yet implemented
			w.WriteHeader(http.StatusNotImplemented)
			break
		}
	}
}

func (app *App) createConnection() (*Request, *Response, *Next) {
	req := &Request{params: make(map[string]string), query: make(map[string]string)}
	res := &Response{end: false, statusCode: 0, body: ""}
	next := &Next{next: false, route: ""}
	return req, res, next
}

func (app *App) Get(path string, callbacks ...Callback) error {
	return app.handler.register(http.MethodGet, path, &UserHandler{app: app, callbacks: callbacks})
}
