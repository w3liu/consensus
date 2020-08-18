package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

type Reader struct {
	io.Reader
	Total   int64
	Current int64
}

func (r *Reader) Read(p []byte) (n int, err error) {
	n, err = r.Reader.Read(p)
	r.Current += int64(n)
	fmt.Printf("\r进度 %.2f%%", float64(r.Current*10000/r.Total)/100)
	return
}

func DownloadFileProgress(url, fileName string) {
	r, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer func() {
		_ = r.Body.Close()
	}()
	f, err := os.Create(fileName)
	if err != nil {
		panic(err)
	}
	defer func() {
		_ = f.Close()
	}()
	reader := &Reader{
		Reader: r.Body,
		Total:  r.ContentLength,
	}

	_, _ = io.Copy(f, reader)
}
