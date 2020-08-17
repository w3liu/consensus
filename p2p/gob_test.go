package p2p

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"github.com/w3liu/consensus/types"
	"testing"
)

func TestEncode(t *testing.T) {
	obj := types.PacketMsg{
		ChannelID: 0x01,
		EOF:       0x01,
		Data:      []byte("hello gob"),
	}

	var buf bytes.Buffer

	enc := gob.NewEncoder(&buf)

	err := enc.Encode(obj)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println("len", len(buf.Bytes()))

	var msg bean.PacketMsg
	var msg types.PacketMsg

	dec := gob.NewDecoder(&buf)

	err = dec.Decode(&msg)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println("data", string(obj.Data))
}
