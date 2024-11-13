package main

import (
	"errors"
	"fmt"

	"github.com/Eandalf/expressgo"
	"github.com/Eandalf/expressgo/bodyparser"
	"github.com/Eandalf/expressgo/cors"
)

func main() {
	config := expressgo.Config{}
	// config.AllowHost = true
	// config.Coarse = true
	app := expressgo.CreateServer(config)

	// app.Set("case sensitive routing", true)

	app.UseGlobal(cors.Use())

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

	app.Get("/test/req/get", func(req *expressgo.Request, res *expressgo.Response, next *expressgo.Next) {
		res.Send(req.Get("Accept"))
	})

	app.Get("/test/req/header", func(req *expressgo.Request, res *expressgo.Response, next *expressgo.Next) {
		res.Send(req.Header("Accept"))
	})

	app.Get("/test/res/set", func(req *expressgo.Request, res *expressgo.Response, next *expressgo.Next) {
		res.Set("Access-Control-Allow-Origin", "example.com")
		res.Set("Access-Control-Allow-Origin", "*")
	})

	app.Get("/test/res/get", func(req *expressgo.Request, res *expressgo.Response, next *expressgo.Next) {
		res.Set("Access-Control-Allow-Origin", "*")
		res.Send(res.Get("Access-Control-Allow-Origin"))
	})

	app.Get("/test/res/append", func(req *expressgo.Request, res *expressgo.Response, next *expressgo.Next) {
		res.Append("Access-Control-Allow-Origin", "example.com")
		res.Append("Access-Control-Allow-Origin", "google.com")
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

	app.Get("/test/error", func(req *expressgo.Request, res *expressgo.Response, next *expressgo.Next) {
		next.Err = errors.New("raised error in /test/error")
	})

	app.Post("/test/body/base", bodyparser.Json(), func(req *expressgo.Request, res *expressgo.Response, next *expressgo.Next) {
		if j, ok := req.Body.(expressgo.BodyJsonBase); ok {
			if t, ok := j["test"]; ok {
				if s, ok := t.(string); ok {
					res.Send(s)
				}
			}
		}

		res.Send("body parsing failed")
	})

	type Test struct {
		Test string `json:"test"`
	}

	app.Post("/test/body/type", bodyparser.Json(bodyparser.JsonConfig{Receiver: &Test{}}), func(req *expressgo.Request, res *expressgo.Response, next *expressgo.Next) {
		if t, ok := req.Body.(*Test); ok {
			res.Send(t.Test)
		}

		res.Send("body parsing failed")
	})

	app.Post("/test/body/raw", bodyparser.Raw(), func(req *expressgo.Request, res *expressgo.Response, next *expressgo.Next) {
		if b, ok := req.Body.([]byte); ok {
			s := string(b)
			res.Send(s)
		}
	})

	app.UseError(
		"/test/error",
		func(err error, req *expressgo.Request, res *expressgo.Response, next *expressgo.Next) {
			req.Params["error0"] = err.Error()
			next.Err = errors.New("raised error in /test/error 1st error handler")
		}, func(err error, req *expressgo.Request, res *expressgo.Response, next *expressgo.Next) {
			req.Params["error1"] = err.Error()
			next.Err = errors.New("raised error in /test/error 2nd error handler")
		},
	)

	app.UseGlobalError(func(err error, req *expressgo.Request, res *expressgo.Response, next *expressgo.Next) {
		output := ""
		for k, v := range req.Params {
			output += fmt.Sprintf("%s: %s<br />", k, v)
		}
		res.Send(output + "global error")
	})

	app.Listen(8080)
}
