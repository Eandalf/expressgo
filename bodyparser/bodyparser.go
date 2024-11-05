package bodyparser

import (
	"strings"

	"github.com/Eandalf/expressgo"
)

// Check if the received type is the expected type.
func isContentType(value string, expectedType string) bool {
	receivedType := strings.TrimSpace(strings.Split(value, ";")[0])
	return expectedType == receivedType
}

func Json(jsonConfig ...JsonConfig) expressgo.Callback {
	return createJsonParser(jsonConfig)
}
