package main

import (
	"github.com/Eandalf/expressgo"
)

func main() {
	app := expressgo.CreateServer()
	app.Listen(-1)
}
