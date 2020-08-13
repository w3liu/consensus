package gobio

import (
	"bufio"
	"bytes"
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
		buf:    &buffer,
		closer: closer,
	}
}

type gobWriter struct {
	w      io.Writer
	buf    *bytes.Buffer
	closer io.Closer
}

func (w *gobWriter) WriteMsg(msg bean.Message) (int, error) {
	enc := gob.NewEncoder(w.buf)
	if err := enc.Encode(msg); err != nil {
		return 0, err
	}
	n, err := w.w.Write(w.buf.Bytes())

	if w, ok := w.w.(*bufio.Writer); ok {
		err := w.Flush()
		if err != nil {
			return 0, err
		}
	}

	defer w.buf.Reset()
	return n, err
}

func (w *gobWriter) Close() error {
	if w.closer != nil {
		return w.closer.Close()
	}
	return nil
}
