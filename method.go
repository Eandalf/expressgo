package expressgo

import (
	"net/http"
)

const (
	configKeyCaseSensitive = "case sensitive routing"
)

var allMethods = [...]string{
	http.MethodGet,
	http.MethodHead,
	http.MethodPost,
	http.MethodPut,
	http.MethodPatch,
	http.MethodDelete,
	http.MethodConnect,
	http.MethodOptions,
	http.MethodTrace,
}

// Provided with an app global data table.
//
// Set data: app.Set(string, interface{})
//
// Get data: app.GetData(string) interface{}
//
// Other than setting data into the app global data table, the method could set app configuration options.
//
// e.g., app.Set("case sensitive routing", true)
func (app *App) Set(key string, value interface{}) {
	switch key {
	case configKeyCaseSensitive:
		if isCaseSensitive, ok := value.(bool); ok {
			app.config.caseSensitive = isCaseSensitive
		}
	}

	app.data[key] = value
}

// Get data from the app global data table.
func (app *App) GetData(key string) interface{} {
	if data, ok := app.data[key]; ok {
		return data
	}
	return nil
}

// To mount callbacks as middlewares to the path with all http methods.
//
// The order of declaration matters.
func (app *App) Use(path string, callbacks ...Callback) error {
	for _, method := range allMethods {
		err := app.handler.register(method, path, &UserHandler{app: app, callbacks: callbacks})
		if err != nil {
			return err
		}
	}
	return nil
}

// To catch all http verbs on a path.
//
// Although the implementation is basically the same as app.Use, app.Use is for middlewares, app.All is for http verbs.
//
// It is more semantically correct to use app.All for all http verbs.
func (app *App) All(path string, callbacks ...Callback) error {
	return app.Use(path, callbacks...)
}

func (app *App) Get(path string, callbacks ...Callback) error {
	return app.handler.register(http.MethodGet, path, &UserHandler{app: app, callbacks: callbacks})
}
