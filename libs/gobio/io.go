package gobio

import (
	"github.com/w3liu/consensus/types"
	"io"
)

type Writer interface {
	WriteMsg(types.Message) (int, error)
}

type WriteCloser interface {
	Writer
	io.Closer
}

type Reader interface {
	ReadMsg(types.Message) error
}

type ReadCloser interface {
	Reader
	io.Closer
}

type byteReader struct {
	io.Reader
	bytes []byte
}

func newByteReader(r io.Reader) *byteReader {
	return &byteReader{
		Reader: r,
		bytes:  make([]byte, 1),
	}
}

func (r *byteReader) ReadByte() (byte, error) {
	_, err := r.Read(r.bytes)
	if err != nil {
		return 0, err
	}
	return r.bytes[0], nil
}
