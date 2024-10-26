package bodyparser

import (
	"encoding/json"

	"github.com/Eandalf/expressgo"
)

type JsonConfig struct {
	Receiver interface{}
}

func createJsonParser(jsonConfig []JsonConfig) expressgo.Callback {
	var parser func(*expressgo.Request, *expressgo.Response, *expressgo.Next)
	if len(jsonConfig) > 0 {
		parser = func(req *expressgo.Request, res *expressgo.Response, next *expressgo.Next) {
			// only intercept the request body if Content-Type is set to application/json
			if req.Native.Header.Get("Content-Type") == "application/json" {
				body := jsonConfig[0].Receiver
				err := json.NewDecoder(req.Native.Body).Decode(body)
				if err != nil {
					next.Err = err
				} else {
					req.Body = body
				}
			}

			// proceed to the next callback
			next.Next = true
			next.Route = true
		}
	} else {
		parser = func(req *expressgo.Request, res *expressgo.Response, next *expressgo.Next) {
			// only intercept the request body if Content-Type is set to application/json
			if req.Native.Header.Get("Content-Type") == "application/json" {
				var body expressgo.BodyJsonBase
				err := json.NewDecoder(req.Native.Body).Decode(&body)
				if err != nil {
					next.Err = err
				} else {
					req.Body = body
				}
			}

			// proceed to the next callback
			next.Next = true
			next.Route = true
		}
	}

	return parser
}
