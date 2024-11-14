package bodyparser

import (
	"errors"
	"io"
	"strings"

	"github.com/Eandalf/expressgo"
)

var ErrCu = errors.New("415: charset.unsupported")
var ErrEtl = errors.New("413: entity.too.large")
var ErrEvf = errors.New("403: entity.verify.failed")

type reader struct {
	// underlying reader
	R io.Reader
	// max bytes remaining
	N        int64
	req      *expressgo.Request
	res      *expressgo.Response
	encoding string
	// verify the stream
	verify  Verify
	checked bool
}

func (r *reader) Read(p []byte) (n int, err error) {
	// init, only run once
	if !r.checked {
		if !strings.HasPrefix(r.encoding, "utf-") {
			return 0, ErrCu
		}

		r.checked = true
	}

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
