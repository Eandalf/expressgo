package bodyparser

import (
	"io"
	"regexp"
	"strings"

	"github.com/Eandalf/expressgo"
)

// Check if the received type is the expected type.
func isContentType(value string, expectedType any) bool {
	receivedType := strings.ToLower(strings.TrimSpace(strings.Split(value, ";")[0]))

	if e, ok := expectedType.(string); ok {
		return mimeMatch(receivedType, normalize(e))
	} else if es, ok := expectedType.([]string); ok {
		for _, e := range es {
			if mimeMatch(receivedType, normalize(e)) {
				return true
			}
		}
	}

	return false
}

var charsetMatch = regexp.MustCompile(`charset=([-\w]+)`)

// get charset from http header content-type
func getCharset(value string) string {
	values := strings.Split(value, ";")
	for _, v := range values {
		matches := charsetMatch.FindStringSubmatch(v)
		if len(matches) == 2 {
			return strings.ToLower(matches[1])
		}
	}

	// default to utf-8
	return "utf-8"
}

type Verify func(*expressgo.Request, *expressgo.Response, []byte, string) error

type readOption struct {
	limit    int64
	req      *expressgo.Request
	res      *expressgo.Response
	encoding string
	verify   Verify
}

// read the body stream
func read(r io.Reader, option *readOption) io.Reader {
	return &reader{
		r,
		option.limit,
		option.req,
		option.res,
		option.encoding,
		option.verify,
		false,
	}
}

func Json(jsonConfig ...JsonConfig) expressgo.Callback {
	return createJsonParser(jsonConfig)
}

func Raw(rawConfig ...RawConfig) expressgo.Callback {
	return createRawParser(rawConfig)
}

func Text(textConfig ...TextConfig) expressgo.Callback {
	return createTextParser(textConfig)
}
