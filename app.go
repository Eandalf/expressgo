package expressgo

import (
	"log"
	"net/http"
	"strconv"
)

type appConfig struct {
	allowHost     bool
	coarse        bool
	caseSensitive bool
}

type App struct {
	config *appConfig
	// global data table for app
	data    map[string]interface{}
	handler *Handler
	// multiple lists of callbacks associated with a route, routeA -> [[c11, c12, c13], [c21, c22]]
	callbacks map[string][][]Callback
	// params associated with a route, routeA -> [[param1, param2], [param3]]
	params map[string][][]string
}

type Config struct {
	AllowHost bool
	Coarse    bool
}

func CreateServer(config ...Config) App {
	mux := http.NewServeMux()

	// perform the configuration, config is made to a slice to mimic behaviors of optional parameters
	app := App{
		config:    &appConfig{},
		data:      map[string]interface{}{},
		handler:   &Handler{mux: mux},
		callbacks: map[string][][]Callback{},
		params:    map[string][][]string{},
	}
	app.handler.app = &app
	if len(config) > 0 {
		c := config[0]
		app.config.allowHost = c.AllowHost
		app.config.coarse = c.Coarse
	}

	return app
}

func (app *App) Listen(port int) {
	log.Println("expressgo listens to port: " + strconv.Itoa(port))
	err := http.ListenAndServe(":"+strconv.Itoa(port), app.handler)
	if err != nil {
		log.Fatalln(err)
	}
}
