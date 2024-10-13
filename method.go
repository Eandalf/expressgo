package expressgo

import (
	"net/http"
)

func (app *App) createConnection() (*Request, *Response, *Next) {
	req := &Request{params: make(map[string]string), query: make(map[string]string)}
	res := &Response{end: false, statusCode: 0, body: ""}
	next := &Next{next: false, route: ""}
	return req, res, next
}

func (app *App) Get(path string, callbacks ...Callback) error {
	return app.handler.register(http.MethodGet, path, &UserHandler{app: app, callbacks: callbacks})
}
