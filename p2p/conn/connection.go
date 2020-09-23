package conn

import (
	"bufio"
	"fmt"
	"github.com/w3liu/consensus/libs/gobio"
	"github.com/w3liu/consensus/libs/timer"
	"github.com/w3liu/consensus/log"
	"github.com/w3liu/consensus/types"
	"go.uber.org/zap"
	"io"
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
	flushTimer      *timer.ThrottleTimer
}

type receiveCbFunc func(chID byte, msgBytes []byte)

type errorCbFunc func(interface{})

func NewMConnection(
	conn net.Conn,
	chDescs []*ChannelDescriptor,
	onReceive receiveCbFunc,
	onError errorCbFunc,
) *MConnection {
	mconn := &MConnection{
		conn:          conn,
		bufConnReader: bufio.NewReaderSize(conn, minReadBufferSize),
		bufConnWriter: bufio.NewWriterSize(conn, minWriteBufferSize),
		send:          make(chan struct{}, 1),
		pong:          make(chan struct{}, 1),
		onReceive:     onReceive,
		onError:       onError,
	}
	var channelsIdx = map[byte]*Channel{}
	var channels = make([]*Channel, 0)

	for _, desc := range chDescs {
		channel := newChannel(mconn, *desc)
		channelsIdx[channel.desc.ID] = channel
		channels = append(channels, channel)
	}
	mconn.channels = channels
	mconn.channelsIdx = channelsIdx

	return mconn
}

func (c *MConnection) OnStart() error {
	c.flushTimer = timer.NewThrottleTimer("flush", defaultFlushThrottle)
	c.quitSendRoutine = make(chan struct{})
	c.doneSendRoutine = make(chan struct{})
	c.quitRecvRoutine = make(chan struct{})
	go c.sendRoutine()
	go c.recvRoutine()
	return nil
}

func (c *MConnection) OnStop() {
	if c.stopServices() {
		return
	}
	c.conn.Close()
}

func (c *MConnection) flush() {
	err := c.bufConnWriter.Flush()
	if err != nil {
		log.Error("MConnection flush failed", zap.Error(err))
	}
}

func (c *MConnection) stopServices() (alreadyStopped bool) {
	c.stopMtx.Lock()
	defer c.stopMtx.Unlock()

	select {
	case <-c.quitSendRoutine:
		return true
	default:
	}

	select {
	case <-c.quitRecvRoutine:
		return true
	default:
	}

	// inform the recvRouting that we are shutting down
	close(c.quitRecvRoutine)
	close(c.quitSendRoutine)
	return false
}

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
		log.Info("Send failed", zap.Any("channel", chID), zap.Any("conn", c), zap.Any("msgBytes", fmt.Sprintf("%X", msgBytes)))
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
	var leastChannel *Channel
	for _, channel := range c.channels {
		// If nothing to send, skip this channel
		if !channel.isSendPending() {
			continue
		}
		leastChannel = channel
	}
	// Nothing to send?
	if leastChannel == nil {
		return true
	}
	// c.Logger.Info("Found a msgPacket to send")
	_, err := leastChannel.writePacketMsgTo(c.bufConnWriter)
	if err != nil {
		log.Error("writePacketMsgTo error", zap.Error(err))
		return true
	}
	//c.flush()
	c.flushTimer.Set()
	return false
}

func (c *MConnection) sendRoutine() {
	defer c._recover()
	for {
		select {
		case <-c.flushTimer.Ch:
			// NOTE: flushTimer.Set() must be called every time
			// something is written to .bufConnWriter.
			c.flush()
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
	reader := gobio.NewReader(c.conn)
FOR_LOOP:
	for {
		var packet types.PacketMsg
		err := reader.ReadMsg(&packet)
		if err != nil {
			select {
			case <-c.quitRecvRoutine:
				break FOR_LOOP
			default:
			}
			if err == io.EOF {
				log.Info("Connection is closed @ recvRoutine (likely by the other side)", zap.Any("conn", c))
			} else {
				log.Error("Connection failed @ recvRoutine (reading byte)", zap.Any("conn", c), zap.Error(err))
			}
			c.onError(err)
			break FOR_LOOP
		}

		channel, ok := c.channelsIdx[byte(packet.ChannelID)]
		if !ok || channel == nil {
			err := fmt.Errorf("unknown channel %X", packet.ChannelID)
			log.Error("Connection failed2 @ recvRoutine", zap.Any("conn", c), zap.Error(err))
			break FOR_LOOP
		}

		msgBytes, err := channel.recvPacketMsg(packet)
		if err != nil {
			log.Info("Connection failed3 @ recvRoutine", zap.Any("conn", c), zap.Error(err))
			break FOR_LOOP
		}
		if msgBytes != nil {
			// NOTE: This means the reactor.Receive runs in the same thread as the p2p recv routine
			c.onReceive(byte(packet.ChannelID), msgBytes)
		}
	}
}

func (c *MConnection) _recover() {
	if r := recover(); r != nil {
		log.Error("MConnection panicked", zap.Any("err", r), zap.Any("stack", string(debug.Stack())))
	}
}
