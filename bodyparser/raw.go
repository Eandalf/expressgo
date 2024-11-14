package bodyparser

import (
	"io"

	"github.com/Eandalf/expressgo"
)

type RawConfig struct {
	Type     any
	Limit    any
	limitNum int64
	Verify   Verify
}

func createRawParser(rawConfig []RawConfig) expressgo.Callback {
	config := RawConfig{
		Type:  "application/octet-stream",
		Limit: "100kb",
	}

	if len(rawConfig) > 0 {
		userConfig := rawConfig[0]

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
				req.Body = body
			}
		}

		// proceed to the next callback
		next.Next = true
		next.Route = true
	}

	return parser
}
