package p2p

import (
	"github.com/w3liu/consensus/p2p/conn"
	"net"
)

type peerConn struct {
	outbound   bool
	conn       net.Conn
	socketAddr *NetAddress
	ip         net.IP
}

func newPeerConn(outbound bool, conn net.Conn, socketAddr *NetAddress) peerConn {
	return peerConn{
		outbound:   outbound,
		conn:       conn,
		socketAddr: socketAddr,
	}
}

type peer struct {
	peerConn
	mconn    *conn.MConnection
	nodeInfo NodeInfo
	channels []byte
}

func newPeer(pc peerConn, nodeInfo NodeInfo) *peer {
	return &peer{
		peerConn: pc,
		nodeInfo: nodeInfo,
	}
}

func (p *peer) Send(chID byte, msgBytes []byte) bool {
	if !p.hasChannel(chID) {
		return false
	}
	res := p.mconn.Send(chID, msgBytes)
	return res
}

func (p *peer) hasChannel(chID byte) bool {
	for _, ch := range p.channels {
		if ch == chID {
			return true
		}
	}
	return false
}
