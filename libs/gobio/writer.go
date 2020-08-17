package gobio

import (
	"bytes"
	"encoding/binary"
	"encoding/gob"
	"github.com/w3liu/consensus/bean"
	"io"
)

func NewWriter(w io.Writer) WriteCloser {
	var closer io.Closer
	var buffer bytes.Buffer
	if c, ok := w.(io.Closer); ok {
		closer = c
	}
	return &gobWriter{
		w:      w,
		buffer: &buffer,
		lenBuf: make([]byte, binary.MaxVarintLen64),
		closer: closer,
	}
}

type gobWriter struct {
	w      io.Writer
	buffer *bytes.Buffer
	lenBuf []byte
	closer io.Closer
}

func (w *gobWriter) WriteMsg(msg bean.Message) (int, error) {
	enc := gob.NewEncoder(w.buffer)
	if err := enc.Encode(msg); err != nil {
		return 0, err
	}
	length := uint64(len(w.buffer.Bytes()))
	n := binary.PutUvarint(w.lenBuf, length)
	_, err := w.w.Write(w.lenBuf[:n])
	if err != nil {
		return 0, err
	}
	n, err = w.w.Write(w.buffer.Bytes())
	return n, err
}

func (w *gobWriter) Close() error {
	if w.closer != nil {
		return w.closer.Close()
	}
	return nil
}
