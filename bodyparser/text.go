package bodyparser

import (
	"io"

	"github.com/Eandalf/expressgo"
)

type TextConfig struct{}

func createTextParser(textConfig []TextConfig) expressgo.Callback {
	var parser expressgo.Callback

	if len(textConfig) > 0 {
		parser = func(req *expressgo.Request, res *expressgo.Response, next *expressgo.Next) {
			if isContentType(req.Native.Header.Get("Content-Type"), "text/plain") {
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
	} else {
		parser = func(req *expressgo.Request, res *expressgo.Response, next *expressgo.Next) {
			if isContentType(req.Native.Header.Get("Content-Type"), "text/plain") {
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
	}

	return parser
}
