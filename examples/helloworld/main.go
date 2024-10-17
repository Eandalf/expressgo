package main

import (
	"github.com/Eandalf/expressgo"
)

func main() {
	config := expressgo.Config{}
	// config.AllowHost = true
	// config.Coarse = true
	app := expressgo.CreateServer(config)

	// app.Use("case sensitive routing", true)

	app.Get("/test/status", func(req *expressgo.Request, res *expressgo.Response, next *expressgo.Next) {
		res.SendStatus(201)
	})

	app.Get("/test/body", func(req *expressgo.Request, res *expressgo.Response, next *expressgo.Next) {
		res.Status(202).Send("Hello from /test/body")
	})

	app.Get("/test/end", func(req *expressgo.Request, res *expressgo.Response, next *expressgo.Next) {
		res.End()
	})

	app.Get("127.0.0.1/test/host", func(req *expressgo.Request, res *expressgo.Response, next *expressgo.Next) {
		res.Send("Hello from 127.0.0.1/test/host")
	})

	app.Get("/", func(req *expressgo.Request, res *expressgo.Response, next *expressgo.Next) {
		res.Send("Hello from root")
	})

	app.Get(
		"/test/next",
		func(req *expressgo.Request, res *expressgo.Response, next *expressgo.Next) {
			req.Params["id"] = "101"
			next.Next = true
		},
		func(req *expressgo.Request, res *expressgo.Response, next *expressgo.Next) {
			res.Send("id: " + req.Params["id"])
		},
	)

	app.Get("/test/next/route", func(req *expressgo.Request, res *expressgo.Response, next *expressgo.Next) {
		req.Params["id"] = "101"
		next.Route = true
	})

	app.Get("/test/next/route", func(req *expressgo.Request, res *expressgo.Response, next *expressgo.Next) {
		res.Send("id: " + req.Params["id"])
	})

	app.Listen(8080)
}
