package bodyparser

import (
	"io"

	"github.com/Eandalf/expressgo"
)

type TextConfig struct {
	Type           any
	Inflate        bool
	Limit          any
	limitNum       int64
	Verify         Verify
	DefaultCharset string
}

func createTextParser(textConfig []TextConfig) expressgo.Callback {
	config := TextConfig{
		Type:           "text/plain",
		Inflate:        true,
		Limit:          "100kb",
		DefaultCharset: "utf-8",
	}

	if len(textConfig) > 0 {
		userConfig := textConfig[0]

		if userConfig.Type != nil {
			config.Type = userConfig.Type
		}
		if !userConfig.Inflate {
			config.Inflate = userConfig.Inflate
		}
		if userConfig.Limit != nil {
			config.Limit = userConfig.Limit
		}
		if userConfig.Verify != nil {
			config.Verify = userConfig.Verify
		}
		if userConfig.DefaultCharset != "" {
			config.DefaultCharset = userConfig.DefaultCharset
		}
	}

	config.limitNum = parseByte(config.Limit)

	parser := func(req *expressgo.Request, res *expressgo.Response, next *expressgo.Next) {
		if isContentType(req.Native.Header.Get("Content-Type"), config.Type) {
			charset := getCharset(req.Native.Header.Get("Content-Type"), config.DefaultCharset)

			stream, sErr := getStream(
				req.Native.Body,
				&readOption{
					config.Inflate,
					config.limitNum,
					req,
					res,
					config.Verify,
				},
				req.Native.Header.Get("Content-Encoding"),
				charset,
			)
			if sErr != nil {
				next.Err = sErr
				return
			}

			body, err := io.ReadAll(stream)

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
