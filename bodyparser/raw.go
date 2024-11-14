package bodyparser

import (
	"io"

	"github.com/Eandalf/expressgo"
)

type RawConfig struct {
	Type any
}

func createRawParser(rawConfig []RawConfig) expressgo.Callback {
	config := RawConfig{
		Type: "application/octet-stream",
	}

	if len(rawConfig) > 0 {
		userConfig := rawConfig[0]

		if userConfig.Type != nil {
			config.Type = userConfig.Type
		}
	}

	parser := func(req *expressgo.Request, res *expressgo.Response, next *expressgo.Next) {
		if isContentType(req.Native.Header.Get("Content-Type"), config.Type) {
			body, err := io.ReadAll(req.Native.Body)

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
