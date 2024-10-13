package expressgo

import (
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type App struct {
	handler Handler
}

func CreateServer() App {
	mux := http.NewServeMux()
	return App{handler: Handler{mux: mux}}
}

func (app *App) Listen(port int) {
	log.Println("expressgo listens to port: " + strconv.Itoa(port))
	err := http.ListenAndServe(":"+strconv.Itoa(port), &app.handler)
	if err != nil {
		log.Fatalln(err)
	}
}

// For path registration

func (app *App) isHostIncluded(path string) bool {
	return path[0] != '/'
}

func (app *App) makePrecise(path string) string {
	return path + "/{$}"
}

func (app *App) pathToLower(path string) string {
	return strings.ToLower(path)
}

func (app *App) register(method string, path string, handler http.Handler) error {
	if app.isHostIncluded(path) {
		return errors.New("path cannot contain host")
	}

	p := app.makePrecise(app.pathToLower(path))
	if method != "" {
		p = method + " " + p
	}

	app.handler.mux.Handle(p, handler)
	return nil
}

// For processing requests

type Handler struct {
	mux *http.ServeMux
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	r.URL.Path = strings.ToLower(r.URL.Path)
	h.mux.ServeHTTP(w, r)
}
