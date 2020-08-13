package conn

import (
	"fmt"
	"net"
	"testing"
	"time"
)

func createTestMConnection(conn net.Conn) *MConnection {
	onReceive := func(chID byte, msgBytes []byte) {
	}
	onError := func(r interface{}) {
	}
	c := createMConnectionWithCallbacks(conn, onReceive, onError)
	return c
}

func createMConnectionWithCallbacks(
	conn net.Conn,
	onReceive func(chID byte, msgBytes []byte),
	onError func(r interface{}),
) *MConnection {

	chDescs := []*ChannelDescriptor{{ID: 0x01, Priority: 1, SendQueueCapacity: 1}}
	c := NewMConnection(conn, chDescs, onReceive, onError)
	return c
}

func TestMConnectionReceive(t *testing.T) {
	server, client := net.Pipe()
	defer server.Close()
	defer client.Close()

	serverConn := createTestMConnection(server)
	serverConn.OnStart()
	defer serverConn.OnStop()

	receivedCh := make(chan []byte)
	errorsCh := make(chan interface{})
	onReceive := func(chID byte, msgBytes []byte) {
		receivedCh <- msgBytes
	}
	onError := func(r interface{}) {
		errorsCh <- r
	}
	clientConn := createMConnectionWithCallbacks(client, onReceive, onError)
	clientConn.OnStart()
	defer clientConn.OnStop()

	msg := "hello"

	for i := 0; i < 5; i++ {
		txt := fmt.Sprintf("%s_%d", msg, i)
		serverConn.Send(0x01, []byte(txt))
		time.Sleep(time.Second)
	}
	for {
		select {
		case receivedBytes := <-receivedCh:
			fmt.Println("receivedBytes", string(receivedBytes))
		case err := <-errorsCh:
			t.Fatalf("Expected %s, got %+v", msg, err)
		case <-time.After(50 * time.Second):
			t.Fatalf("Did not receive %s message in 50s", msg)
		}
	}

}

func TestPipe(t *testing.T) {
	server, client := net.Pipe()

	for i := 0; i < 10; i++ {
		go func() {
			_, err := server.Write([]byte("hello"))
			if err != nil {
				t.Fatal(err)
			}
		}()
	}
	buf := make([]byte, 100)
	var i = 0
	for {
		n, err := client.Read(buf)
		if err != nil {
			t.Fatal(err)
		}
		fmt.Println(string(buf[:n]))
		i++
		if i >= 9 {
			break
		}
	}
}
