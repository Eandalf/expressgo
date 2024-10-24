package main

import (
	"fmt"

	"github.com/Eandalf/expressgo"
)

func main() {
	config := expressgo.Config{}
	// config.AllowHost = true
	// config.Coarse = true
	app := expressgo.CreateServer(config)

	// app.Set("case sensitive routing", true)

	app.UseGlobal(func(req *expressgo.Request, res *expressgo.Response, next *expressgo.Next) {
		req.Params["global"] = "global"
		next.Route = true
	})

	app.Get("/test/status/1", func(req *expressgo.Request, res *expressgo.Response, next *expressgo.Next) {
		res.SendStatus(201)
	})

	app.Get("/test/status/2", func(req *expressgo.Request, res *expressgo.Response, next *expressgo.Next) {
		res.Status(202).Send("Hello from /test/status/2")
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

	app.Get("/test/params/:one-:two-:three/:four.:five", func(req *expressgo.Request, res *expressgo.Response, next *expressgo.Next) {
		lines := []string{}
		for k, v := range req.Params {
			lines = append(lines, fmt.Sprintf("%s: %s", k, v))
		}

		output := ""
		for _, line := range lines {
			output += line + "<br />"
		}
		res.Send(output)
	})

	app.Get("/test/query", func(req *expressgo.Request, res *expressgo.Response, next *expressgo.Next) {
		output := ""
		for k, v := range req.Query {
			output += fmt.Sprintf("%s: %s<br />", k, v)
		}
		res.Send(output)
	})

	app.Use("/test/use", func(req *expressgo.Request, res *expressgo.Response, next *expressgo.Next) {
		req.Params["id"] = "102"
		next.Route = true
	})

	app.All("/test/use", func(req *expressgo.Request, res *expressgo.Response, next *expressgo.Next) {
		res.Send("id: " + req.Params["id"])
	})

	app.Listen(8080)
}
