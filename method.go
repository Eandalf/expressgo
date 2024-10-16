package expressgo

import (
	"net/http"
)

const (
	useOptionKeyCaseSensitive = "case sensitive routing"
)

func (app *App) Use(key string, value bool) {
	switch key {
	case useOptionKeyCaseSensitive:
		app.caseSensitive = value
	}
}

func (app *App) Get(path string, callbacks ...Callback) error {
	route := http.MethodGet + " " + path

	// register the slice of callbacks with the route formed by the method and the path
	// if the route already exists, push the slice of callbacks to map and not register it to ServeMux
	if _, ok := app.routes[route]; ok {
		existingCallbacks := app.routes[route]
		app.routes[route] = append(existingCallbacks, callbacks)
		return nil
	}

	app.routes[route] = [][]Callback{callbacks}
	return app.handler.register(http.MethodGet, path, &UserHandler{app: app, callbacks: callbacks})
}
