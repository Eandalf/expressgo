package bodyparser

import (
	"compress/flate"
	"compress/gzip"
	"compress/lzw"
	"io"
	"slices"
	"strings"

	"github.com/Eandalf/expressgo"
	"golang.org/x/text/encoding/ianaindex"
	"golang.org/x/text/transform"
)

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
		if tv := strings.TrimSpace(v); tv != "" {
			result = append(result, tv)
		}
	}

	return result
}

// Check if all compressions are "identity" (no compression).
func areAllCompressionsIdentity(compressions []string) bool {
	for _, c := range compressions {
		if c != compressionTypeIdentity {
			return false
		}
	}

	return true
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

// Limit the length of stream reading and perform verify
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
	charset string,
) (pipe io.Reader, err error) {
	// get header infos
	compressions := getCompressions(contentEncoding)
	if !option.inflate && !areAllCompressionsIdentity(compressions) {
		err = ErrEu
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
			pipe = flate.NewReader(pipe)
		case compressionTypeGzip:
			pipe, err = gzip.NewReader(pipe)
			if err != nil {
				return
			}
		case compressionTypeCompress:
			pipe = lzw.NewReader(pipe, lzw.LSB, 8)
		default:
			// compression unsupported, raise the error
			err = ErrEu
			return
		}
	}

	// decode the stream based on the charset
	if charset != "" && charset != "utf-8" {
		encoding, eErr := ianaindex.IANA.Encoding(charset)
		if eErr != nil || encoding == nil {
			// charset unsupported, raise the error
			err = ErrCu
			return
		}
		pipe = transform.NewReader(pipe, encoding.NewDecoder().Transformer)
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
