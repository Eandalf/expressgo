package bodyparser

import (
	"io"

	"github.com/Eandalf/expressgo"
)

type TextConfig struct {
	Type any
}

func createTextParser(textConfig []TextConfig) expressgo.Callback {
	config := TextConfig{
		Type: "text/plain",
	}

	if len(textConfig) > 0 {
		userConfig := textConfig[0]

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
				req.Body = string(body)
			}
		}

		// proceed to the next callback
		next.Next = true
		next.Route = true
	}

	return parser
}
