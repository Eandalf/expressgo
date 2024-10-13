package expressgo

import (
	"log"
	"net/http"
	"strconv"
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
