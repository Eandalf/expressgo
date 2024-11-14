package bodyparser

import (
	"errors"
	"io"
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

type limitedReader struct {
	// underlying reader
	R io.Reader
	// max bytes remaining
	N int64
}

var ErrCtl = errors.New("content too large")

func (l *limitedReader) Read(p []byte) (n int, err error) {
	if l.N <= 0 {
		return 0, io.EOF
	}

	if int64(len(p)) > l.N {
		return 0, ErrCtl
	}

	n, err = l.R.Read(p)
	l.N -= int64(n)
	return
}

// read the body stream
func read(r io.Reader, limit int64) io.Reader {
	return &limitedReader{r, limit}
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
