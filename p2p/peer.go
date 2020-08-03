package p2p

import "net"

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

func (pc peerConn) ID() ID {
	return ""
}

type peer struct {
}
