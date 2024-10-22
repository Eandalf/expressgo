package expressgo

import (
	"net/http"
)

const (
	useOptionKeyCaseSensitive = "case sensitive routing"
)

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
	case useOptionKeyCaseSensitive:
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

func (app *App) Get(path string, callbacks ...Callback) error {
	return app.handler.register(http.MethodGet, path, &UserHandler{app: app, callbacks: callbacks})
}
