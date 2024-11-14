package bodyparser

import (
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

func Json(jsonConfig ...JsonConfig) expressgo.Callback {
	return createJsonParser(jsonConfig)
}

func Raw(rawConfig ...RawConfig) expressgo.Callback {
	return createRawParser(rawConfig)
}

func Text(textConfig ...TextConfig) expressgo.Callback {
	return createTextParser(textConfig)
}
