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
	return app.handler.register(http.MethodGet, path, &UserHandler{app: app, callbacks: callbacks})
}
