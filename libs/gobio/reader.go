package gobio

import (
	"bytes"
	"encoding/binary"
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
		buf:    make([]byte, 0),
		closer: closer,
	}
}

type gobReader struct {
	r      io.Reader
	buf    []byte
	closer io.Closer
}

func (r *gobReader) ReadMsg(msg types.Message) error {
	length64, err := binary.ReadUvarint(newByteReader(r.r))

	n, err := r.r.Read(r.buf)
	if err != nil {
		return err
	}
	length := int(length64)
	if length < 0 {
		return fmt.Errorf("message length is 0")
	}

	if len(r.buf) < length {
		r.buf = make([]byte, length)
	}
	buf := r.buf[:length]
	if _, err := io.ReadFull(r.r, buf); err != nil {
		return err
	}
	dec := gob.NewDecoder(bytes.NewBuffer(buf))
	return dec.Decode(msg)
}

func (r *gobReader) Close() error {
	if r.closer != nil {
		return r.closer.Close()
	}
	return nil
}
