package bodyparser

import (
	"io"
	"net/url"

	"github.com/Eandalf/expressgo"
)

type UrlencodedConfig struct {
	Type           any
	Inflate        bool
	Limit          any
	limitNum       int64
	Verify         Verify
	DefaultCharset string
}

func createUrlencodedParser(urlencodedConfig []UrlencodedConfig) expressgo.Callback {
	config := UrlencodedConfig{
		Type:           "application/x-www-form-urlencoded",
		Inflate:        true,
		Limit:          "100kb",
		DefaultCharset: "utf-8",
	}

	if len(urlencodedConfig) > 0 {
		userConfig := urlencodedConfig[0]

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
				return
			}

			vs, pErr := url.ParseQuery(string(body))
			if pErr != nil {
				next.Err = err
				return
			}

			pf := make(expressgo.BodyFormUrlEncoded)
			for k := range vs {
				// we only accept the first value of a key from the values (vs)
				pf[k] = vs.Get(k)
			}

			req.Body = pf
		}

		// proceed to the next callback
		next.Next = true
		next.Route = true
	}

	return parser
}
