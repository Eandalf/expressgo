package bodyparser

import (
	"errors"
	"regexp"
	"strings"

	"github.com/Eandalf/expressgo"
)

var ErrEu = errors.New("415: encoding.unsupported")
var ErrCu = errors.New("415: charset.unsupported")
var ErrEtl = errors.New("413: entity.too.large")
var ErrEvf = errors.New("403: entity.verify.failed")

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

// Get charset from http header content-type.
func getCharset(value string, defaultCharset string) string {
	values := strings.Split(value, ";")
	for _, v := range values {
		matches := charsetMatch.FindStringSubmatch(v)
		if len(matches) == 2 {
			return strings.ToLower(matches[1])
		}
	}

	return strings.ToLower(defaultCharset)
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

func Urlencoded(urlencodedConfig ...UrlencodedConfig) expressgo.Callback {
	return createUrlencodedParser(urlencodedConfig)
}
