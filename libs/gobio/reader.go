package gobio

import (
	"bytes"
	"encoding/gob"
	"github.com/w3liu/consensus/bean"
	"io"
)

func NewReader(r io.Reader) ReadCloser {
	var closer io.Closer
	var buffer bytes.Buffer
	if c, ok := r.(io.Closer); ok {
		closer = c
	}
	return &gobReader{
		r:      r,
		buf:    &buffer,
		closer: closer,
	}
}

type gobReader struct {
	r      io.Reader
	buf    *bytes.Buffer
	closer io.Closer
}

func (r *gobReader) ReadMsg(msg bean.Message) error {
	dec := gob.NewDecoder(r.r)
	return dec.Decode(msg)
}

func (r *gobReader) Close() error {
	if r.closer != nil {
		return r.closer.Close()
	}
	return nil
}
