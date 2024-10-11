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
	err := http.ListenAndServe(":"+strconv.Itoa(port), nil)
	if err != nil {
		log.Println("expressgo listen to: " + strconv.Itoa(port))
	}
}
