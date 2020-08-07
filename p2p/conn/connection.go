package conn

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"net"
	"runtime/debug"
	"sync"
)

type MConnection struct {
	conn            net.Conn
	bufConnReader   *bufio.Reader
	bufConnWriter   *bufio.Writer
	send            chan struct{}
	pong            chan struct{}
	channels        []*Channel
	channelsIdx     map[byte]*Channel
	onReceive       receiveCbFunc
	onError         errorCbFunc
	errored         uint32
	quitSendRoutine chan struct{}
	doneSendRoutine chan struct{}
	quitRecvRoutine chan struct{}
	stopMtx         sync.Mutex
}

type receiveCbFunc func(chID byte, msgBytes []byte)

type errorCbFunc func(interface{})

func (c *MConnection) Send(chID byte, msgBytes []byte) bool {
	channel, ok := c.channelsIdx[chID]
	if !ok {
		return false
	}
	success := channel.sendBytes(msgBytes)
	if success {
		select {
		case c.send <- struct{}{}:
		default:
		}
	} else {
		log.Println("Send failed", "channel", chID, "conn", c, "msgBytes", fmt.Sprintf("%X", msgBytes))
	}
	return success
}

func (c *MConnection) sendSomePacketMsgs() bool {
	for i := 0; i < numBatchPacketMsgs; i++ {
		if c.sendPacketMsg() {
			return true
		}
	}
	return false
}

func (c *MConnection) sendPacketMsg() bool {
	// Choose a channel to create a PacketMsg from.
	// The chosen channel will be the one whose recentlySent/priority is the least.
	var leastRatio float32 = math.MaxFloat32
	var leastChannel *Channel
	for _, channel := range c.channels {
		// If nothing to send, skip this channel
		if !channel.isSendPending() {
			continue
		}
		// Get ratio, and keep track of lowest ratio.
		ratio := float32(channel.recentlySent) / float32(channel.desc.Priority)
		if ratio < leastRatio {
			leastRatio = ratio
			leastChannel = channel
		}
	}

	// Nothing to send?
	if leastChannel == nil {
		return true
	}
	// c.Logger.Info("Found a msgPacket to send")

	// Make & send a PacketMsg from this channel
	_, err := leastChannel.writePacketMsgTo(c.bufConnWriter)
	if err != nil {
		log.Println(err.Error())
		return true
	}
	return false
}

func (c *MConnection) sendRoutine() {
	defer c._recover()
	for {
		select {
		case <-c.send:
			eof := c.sendSomePacketMsgs()
			if !eof {
				// Keep sendRoutine awake.
				select {
				case c.send <- struct{}{}:
				default:
				}
			}
		}
	}
}

func (c *MConnection) recvRoutine() {
	defer c._recover()

}

func (c *MConnection) _recover() {
	if r := recover(); r != nil {
		log.Println("MConnection panicked", "err", r, "stack", string(debug.Stack()))
	}
}
