package bodyparser

import (
	"io"

	"github.com/Eandalf/expressgo"
)

type TextConfig struct {
	Type     any
	Limit    any
	limitNum int64
	Verify   Verify
	// TODO: defaultCharset
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
		if userConfig.Verify != nil {
			config.Verify = userConfig.Verify
		}
	}

	config.limitNum = parseByte(config.Limit)

	parser := func(req *expressgo.Request, res *expressgo.Response, next *expressgo.Next) {
		if isContentType(req.Native.Header.Get("Content-Type"), config.Type) {
			body, err := io.ReadAll(read(req.Native.Body, &readOption{
				config.limitNum,
				req,
				res,
				getCharset(req.Native.Header.Get("Content-Type")),
				config.Verify,
			}))

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
