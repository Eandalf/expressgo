package expressgo

import (
	"log"
	"net/http"
	"strconv"
)

type App struct{}

func CreateServer() App {
	return App{}
}

func (app *App) Listen(port int) {
	log.Println("expressgo listens to port: " + strconv.Itoa(port))
	err := http.ListenAndServe(":"+strconv.Itoa(port), nil)
	if err != nil {
		log.Fatalln(err)
	}
}
