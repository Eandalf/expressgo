package bodyparser

import (
	"encoding/json"
	"io"
	"reflect"

	"github.com/Eandalf/expressgo"
)

type JsonConfig struct {
	Receiver any
	Type     any
	Limit    any
	limitNum int64
	Verify   Verify
}

func createJsonParser(jsonConfig []JsonConfig) expressgo.Callback {
	config := JsonConfig{
		Receiver: &expressgo.BodyJsonBase{},
		Type:     "application/json",
		Limit:    "100kb",
	}

	if len(jsonConfig) > 0 {
		userConfig := jsonConfig[0]

		if userConfig.Receiver != nil && reflect.ValueOf(userConfig.Receiver).Kind() == reflect.Ptr {
			config.Receiver = userConfig.Receiver
		}
		if userConfig.Type != nil {
			config.Type = userConfig.Type
		}
		if userConfig.Limit != nil {
			config.Limit = userConfig.Limit
		}
		if userConfig.Verify != nil {
			config.Verify = userConfig.Verify
		}
	}

	config.limitNum = parseByte(config.Limit)

	parser := func(req *expressgo.Request, res *expressgo.Response, next *expressgo.Next) {
		// only intercept the request body if Content-Type is set to application/json
		if isContentType(req.Native.Header.Get("Content-Type"), config.Type) {
			body := config.Receiver
			err := json.NewDecoder(read(req.Native.Body, &readOption{
				config.limitNum,
				req,
				res,
				getCharset(req.Native.Header.Get("Content-Type")),
				config.Verify,
			})).Decode(body)

			if err != nil {
				// if EOF is read, either Body is blank or Body has be consumed by parsers before, then no-op
				// otherwise, pass the error to error-handling callbacks
				if err != io.EOF {
					next.Err = err
				}
			} else {
				req.Body = body
			}
		}

		// proceed to the next callback
		next.Next = true
		next.Route = true
	}

	return parser
}
