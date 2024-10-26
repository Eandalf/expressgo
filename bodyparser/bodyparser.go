package bodyparser

import "github.com/Eandalf/expressgo"

func Json(jsonConfig ...JsonConfig) expressgo.Callback {
	return createJsonParser(jsonConfig)
}
