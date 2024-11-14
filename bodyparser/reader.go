package bodyparser

import (
	"compress/gzip"
	"compress/lzw"
	"compress/zlib"
	"errors"
	"io"
	"regexp"
	"slices"
	"strings"

	"github.com/Eandalf/expressgo"
)

var ErrEu = errors.New("415: encoding.unsupported")
var ErrCu = errors.New("415: charset.unsupported")
var ErrEtl = errors.New("413: entity.too.large")
var ErrEvf = errors.New("403: entity.verify.failed")

const (
	compressionTypeIdentity = "identity"
	compressionTypeDeflate  = "deflate"
	compressionTypeGzip     = "gzip"
	compressionTypeCompress = "compress"
)

// Get compressions from http header content-encoding.
func getCompressions(value string) []string {
	result := []string{}
	if value == "" {
		return append(result, compressionTypeIdentity)
	}

	values := strings.Split(value, ",")
	for _, v := range values {
		result = append(result, strings.TrimSpace(v))
	}

	return result
}

// Check if all compressions are "identity" (no compression).
func isAllIdentity(compressions []string) bool {
	for _, c := range compressions {
		if c != compressionTypeIdentity {
			return false
		}
	}

	return true
}

// Check if all listed compression methods are supported.
func areAllCompressionsSupported(compressions []string) bool {
	supportedCompressions := []string{
		compressionTypeIdentity,
		compressionTypeDeflate,
		compressionTypeGzip,
		compressionTypeCompress,
	}

	for _, c := range compressions {
		supported := false

		for _, s := range supportedCompressions {
			if c == s {
				supported = true
				break
			}
		}

		if !supported {
			return false
		}
	}

	return true
}

var charsetMatch = regexp.MustCompile(`charset=([-\w]+)`)

// Get charset from http header content-type.
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

type reader struct {
	// underlying reader
	R io.Reader
	// max bytes remaining
	N        int64
	req      *expressgo.Request
	res      *expressgo.Response
	encoding string
	// verify the stream
	verify Verify
}

func (r *reader) Read(p []byte) (n int, err error) {
	if r.N <= 0 {
		return 0, io.EOF
	}

	if int64(len(p)) > r.N {
		return 0, ErrEtl
	}

	if r.verify != nil {
		if err := r.verify(r.req, r.res, p, r.encoding); err != nil {
			return 0, ErrEvf
		}
	}

	n, err = r.R.Read(p)
	r.N -= int64(n)
	return
}

type readOption struct {
	inflate bool
	limit   int64
	req     *expressgo.Request
	res     *expressgo.Response
	verify  Verify
}

// Set up the pipe to read the body stream.
func getStream(
	r io.Reader,
	option *readOption,
	contentEncoding string,
	contentType string,
) (pipe io.Reader, err error) {
	// get header infos
	compressions := getCompressions(contentEncoding)
	if (!option.inflate && !isAllIdentity(compressions)) || !areAllCompressionsSupported(compressions) {
		err = ErrEu
		return
	}

	charset := getCharset(contentType)
	if !strings.HasPrefix(charset, "utf-") {
		err = ErrCu
		return
	}

	// set up the pipeline to decompress
	pipe = r
	cs := compressions
	slices.Reverse(cs)
	for _, c := range cs {
		switch c {
		case compressionTypeIdentity:
			// do nothing
		case compressionTypeDeflate:
			pipe, err = zlib.NewReader(pipe)
			if err != nil {
				return
			}
		case compressionTypeGzip:
			pipe, err = gzip.NewReader(pipe)
			if err != nil {
				return
			}
		case compressionTypeCompress:
			pipe = lzw.NewReader(pipe, lzw.LSB, 8)
		}
	}

	pipe = &reader{
		pipe,
		option.limit,
		option.req,
		option.res,
		charset,
		option.verify,
	}
	return
}
