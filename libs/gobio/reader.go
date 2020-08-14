package gobio

import (
	"bytes"
	"encoding/gob"
	"github.com/w3liu/consensus/types"
	"io"
)

func NewReader(r io.Reader) ReadCloser {
	var closer io.Closer
	if c, ok := r.(io.Closer); ok {
		closer = c
	}
	return &gobReader{
		r:      r,
		buf:    make([]byte, 1024),
		closer: closer,
	}
}

type gobReader struct {
	r      io.Reader
	buf    []byte
	closer io.Closer
}

func (r *gobReader) ReadMsg(msg types.Message) error {
	n, err := r.r.Read(r.buf)
	if err != nil {
		return err
	}
	dec := gob.NewDecoder(bytes.NewBuffer(r.buf[:n]))
	return dec.Decode(msg)
}

func (r *gobReader) Close() error {
	if r.closer != nil {
		return r.closer.Close()
	}
	return nil
}
