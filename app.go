package expressgo

import (
	"log"
	"net/http"
	"strconv"
)

type App struct {
	handler       *Handler
	routes        map[string][][]Callback
	allowHost     bool
	coarse        bool
	caseSensitive bool
}

type Config struct {
	AllowHost bool
	Coarse    bool
}

func CreateServer(config ...Config) App {
	mux := http.NewServeMux()

	// perform the configuration, config is made to a slice to mimic behaviors of optional parameters
	app := App{handler: &Handler{mux: mux}, routes: make(map[string][][]Callback)}
	app.handler.app = &app
	if len(config) > 0 {
		c := config[0]
		app.allowHost = c.AllowHost
		app.coarse = c.Coarse
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
