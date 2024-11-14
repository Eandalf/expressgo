package bodyparser

import (
	"io"

	"github.com/Eandalf/expressgo"
)

type TextConfig struct {
	Type     any
	Limit    any
	limitNum int64
}

func createTextParser(textConfig []TextConfig) expressgo.Callback {
	config := TextConfig{
		Type:  "text/plain",
		Limit: "100kb",
	}

	if len(textConfig) > 0 {
		userConfig := textConfig[0]

		if userConfig.Type != nil {
			config.Type = userConfig.Type
		}
		if userConfig.Limit != nil {
			config.Limit = userConfig.Limit
		}
	}

	config.limitNum = parseByte(config.Limit)

	parser := func(req *expressgo.Request, res *expressgo.Response, next *expressgo.Next) {
		if isContentType(req.Native.Header.Get("Content-Type"), config.Type) {
			body, err := io.ReadAll(read(req.Native.Body, config.limitNum))

			if err != nil {
				next.Err = err
			} else {
				req.Body = string(body)
			}
		}

		// proceed to the next callback
		next.Next = true
		next.Route = true
	}

	return parser
}
