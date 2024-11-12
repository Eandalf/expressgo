package expressgo

import (
	"net/http"
)

const (
	configKeyAppEnv        = "APP_ENV"
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

// Wrap normal callbacks with error checking logics.
func (app *App) wrapCallbacks(callbacks []Callback) []Callback {
	wrappedCallbacks := []Callback{}
	for _, c := range callbacks {
		var wc Callback = func(req *Request, res *Response, next *Next) {
			// if an error needs to be handled, skip this callback
			if req.err != nil {
				return
			}

			c(req, res, next)
		}
		wrappedCallbacks = append(wrappedCallbacks, wc)
	}

	return wrappedCallbacks
}

// Mount the callbacks to all existing routes and future routes.
//
// This is an internal function that does not take wrapping callbacks into consideration.
//
// Callbacks passed into this function would not be wrapped with error-handling logics.
func (app *App) useGlobal(callbacks []Callback) {
	// add global middlewares to all existing routes
	for route := range app.callbacks {
		app.callbacks[route] = append(app.callbacks[route], callbacks)
	}

	// push global middlewares to globalCallbacks for Handler.register to check and push to callbacks
	*app.globalCallbacks = append(*app.globalCallbacks, callbacks)
}

// To mount middlewares to all existing routes and future routes made by app.[Method].
func (app *App) UseGlobal(callbacks ...Callback) {
	wc := app.wrapCallbacks(callbacks)
	app.useGlobal(wc)
}

// Mount callbacks to the path with all http methods
//
// This is an internal function that does not take wrapping callbacks into the consideration.
//
// Callbacks passed into this function would not be wrapped with error-handling logics.
func (app *App) use(path string, callbacks []Callback) error {
	for _, method := range allMethods {
		err := app.handler.register(method, path, &UserHandler{app: app, callbacks: callbacks})
		if err != nil {
			return err
		}
	}
	return nil
}

// To mount callbacks as middlewares to the path with all http methods.
//
// The order of invocation matters.
func (app *App) Use(path string, callbacks ...Callback) error {
	wc := app.wrapCallbacks(callbacks)
	return app.use(path, wc)
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
	wc := app.wrapCallbacks(callbacks)
	return app.handler.register(http.MethodGet, path, &UserHandler{app: app, callbacks: wc})
}

func (app *App) Head(path string, callbacks ...Callback) error {
	wc := app.wrapCallbacks(callbacks)
	return app.handler.register(http.MethodHead, path, &UserHandler{app: app, callbacks: wc})
}

func (app *App) Post(path string, callbacks ...Callback) error {
	wc := app.wrapCallbacks(callbacks)
	return app.handler.register(http.MethodPost, path, &UserHandler{app: app, callbacks: wc})
}

func (app *App) Put(path string, callbacks ...Callback) error {
	wc := app.wrapCallbacks(callbacks)
	return app.handler.register(http.MethodPut, path, &UserHandler{app: app, callbacks: wc})
}

func (app *App) Patch(path string, callbacks ...Callback) error {
	wc := app.wrapCallbacks(callbacks)
	return app.handler.register(http.MethodPatch, path, &UserHandler{app: app, callbacks: wc})
}

func (app *App) Delete(path string, callbacks ...Callback) error {
	wc := app.wrapCallbacks(callbacks)
	return app.handler.register(http.MethodDelete, path, &UserHandler{app: app, callbacks: wc})
}

func (app *App) Connect(path string, callbacks ...Callback) error {
	wc := app.wrapCallbacks(callbacks)
	return app.handler.register(http.MethodConnect, path, &UserHandler{app: app, callbacks: wc})
}

func (app *App) Options(path string, callbacks ...Callback) error {
	wc := app.wrapCallbacks(callbacks)
	return app.handler.register(http.MethodOptions, path, &UserHandler{app: app, callbacks: wc})
}

func (app *App) Trace(path string, callbacks ...Callback) error {
	wc := app.wrapCallbacks(callbacks)
	return app.handler.register(http.MethodTrace, path, &UserHandler{app: app, callbacks: wc})
}

// Wrap error callbacks into callbacks.
func (app *App) wrapErrorCallbacks(errorCallbacks []ErrorCallback) []Callback {
	callbacks := []Callback{}
	for _, ec := range errorCallbacks {
		var c Callback = func(req *Request, res *Response, next *Next) {
			// if no error needs to be handled
			if req.err == nil {
				return
			}

			ec(req.err, req, res, next)

			// the error is consumed
			req.err = nil
		}
		callbacks = append(callbacks, c)
	}

	return callbacks
}

// To mount error handlers on a path with all http methods.
func (app *App) UseError(path string, errorCallbacks ...ErrorCallback) {
	callbacks := app.wrapErrorCallbacks(errorCallbacks)
	app.use(path, callbacks)
}

// To mount error handlers to all routes.
func (app *App) UseGlobalError(errorCallbacks ...ErrorCallback) {
	callbacks := app.wrapErrorCallbacks(errorCallbacks)
	app.useGlobal(callbacks)
}
