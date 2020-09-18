package p2p

import (
	"bufio"
	"fmt"
	"net"
	"testing"
	"time"
)

func TestSelect(t *testing.T) {
	ch := make(chan struct{})
	select {
	case <-ch:
		fmt.Println("send success")
		return
	case <-time.After(time.Second * 5):
		fmt.Println("time over")
		return
	}
}

func TestDial(t *testing.T) {
	// t.Parallel()
	conn, err := net.DialTimeout("tcp", "127.0.0.1:8000", time.Second*10)
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()
	msg := []byte("hello world")

	buffWriter := bufio.NewWriterSize(conn, 1024)

	ticker := time.NewTicker(time.Second * 5)
	for {
		select {
		case <-ticker.C:
			_, err = buffWriter.Write(msg)
			if err != nil {
				t.Fatal(err)
			}
		}
	}
}

func TestListen(t *testing.T) {
	ln, err := net.Listen("tcp", "127.0.0.1:8000")
	if err != nil {
		t.Fatal(err)
	}
	for {
		c, err := ln.Accept()
		if err != nil {
			t.Fatal(err)
		}
		go func(c net.Conn) {
			bufReader := bufio.NewReaderSize(c, 1024)
			buf := make([]byte, 256)
			for {
				n, err := bufReader.Read(buf)
				if err != nil {
					t.Fatal(err)
				}
				fmt.Println(string(buf[:n]))
			}
			c.RemoteAddr()
		}(c)
	}
}
