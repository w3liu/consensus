package pbft

import (
	"fmt"
	"github.com/w3liu/consensus/config"
	"github.com/w3liu/consensus/libs/gobio"
	"github.com/w3liu/consensus/p2p/conn"
	"github.com/w3liu/consensus/types"
	"log"
	"net"
	"time"
)

type State struct {
	address      *types.Address
	seeds        []*types.Address
	mconn        map[string]*conn.MConnection
	ticker       *time.Ticker // 用于定时出块
	lastBlock    Block
	currentBlock Block
	step         int // 0 pending, 1 propose, 2 prevote, 3 precommit, 4 commit
}

func NewState(cfg *config.Config) *State {
	address, err := types.NewAddress(cfg.Peer.Address)
	if err != nil {
		panic(err)
	}
	seeds := make([]*types.Address, 0)
	for _, seed := range cfg.Peer.Seeds {
		addr, err := types.NewAddress(seed)
		if err != nil {
			panic(err)
		}
		if addr.Id == address.Id {
			continue
		}
		seeds = append(seeds, addr)
	}
	return &State{
		address: address,
		seeds:   seeds,
		mconn:   make(map[string]*conn.MConnection),
		ticker:  time.NewTicker(time.Second * 10),
	}
}

func (s *State) Start() error {
	for {
		select {
		case <-s.ticker.C:
			s.Propose()
		}
	}
}
func (s *State) Stop() {

}

func (s *State) accept() {
	ln, err := net.Listen("tcp", s.address.ToIpPortString())
	if err != nil {
		panic(err)
	}
	for {
		c, err := ln.Accept()
		if err != nil {
			fmt.Println("ln.Accept() error", err)
			continue
		}
		go func(c net.Conn) {
			server, err := s.wrapConn(c)
			if err != nil {
				fmt.Println("s.wrapConn(c) error", err)
				return
			}
			err = server.OnStart()
			if err != nil {
				fmt.Println("server.OnStart()", err)
				return
			}

		}(c)
	}
}

func (s *State) dialPeers() {
	for _, addr := range s.seeds {
		if _, ok := s.mconn[addr.Id]; !ok {
			c, err := s.dial(addr)
			if err != nil {
				fmt.Println("s.dial(addr) error", err)
				continue
			}
			c.OnStart()
		}
	}
}

func (s *State) dial(addr *types.Address) (*conn.MConnection, error) {
	c, err := net.DialTimeout("tcp", addr.ToIpPortString(), time.Second*3)
	if err != nil {
		return nil, err
	}
	return s.wrapConn(c)
}

func (s *State) wrapConn(c net.Conn) (*conn.MConnection, error) {
	peer, err := s.handshake(c, time.Second*3, &types.NodeInfoMessage{
		Id:   s.address.Id,
		Ip:   s.address.Ip.String(),
		Port: s.address.Port,
	})
	if err != nil {
		fmt.Println("handshake error", err)
		c.Close()
		return nil, err
	}
	if _, ok := s.mconn[peer.Id]; ok {
		fmt.Println("connection is existed")
		c.Close()
		return nil, err
	}

	onError := func(r interface{}) {
		// 移除节点
		if _, ok := s.mconn[peer.Id]; ok {
			delete(s.mconn, peer.Id)
		}
		fmt.Println("onError", r)
	}
	chDescs := []*conn.ChannelDescriptor{{ID: 0x01, Priority: 1, SendQueueCapacity: 1}}
	server := conn.NewMConnection(c, chDescs, s.ReceiveMsg, onError)
	return server, nil
}

func (s *State) ReceiveMsg(chID byte, msgBytes []byte) {

}

func (s *State) CreateConn() {

}

func (s *State) handshake(c net.Conn, timeout time.Duration, nodeInfo *types.NodeInfoMessage) (*types.NodeInfoMessage, error) {
	if err := c.SetDeadline(time.Now().Add(timeout)); err != nil {
		return nil, err
	}
	var (
		errc         = make(chan error, 2)
		peerNodeInfo = &types.NodeInfoMessage{}
		ourNodeInfo  = nodeInfo
	)
	go func(errc chan<- error, c net.Conn) {
		_, err := gobio.NewWriter(c).WriteMsg(ourNodeInfo)
		errc <- err
	}(errc, c)

	go func(errc chan<- error, c net.Conn) {
		err := gobio.NewReader(c).ReadMsg(peerNodeInfo)
		errc <- err
	}(errc, c)

	for i := 0; i < cap(errc); i++ {
		err := <-errc
		if err != nil {
			return nil, err
		}
	}
	return peerNodeInfo, c.SetDeadline(time.Time{})
}

func (s *State) Propose() {
	if len(s.mconn) < 2 {
		log.Println("len(s.mconn) < 2")
		return
	}
	if s.step == 0 && s.address.Id == "A" {
		s.currentBlock = Block{
			Height: s.lastBlock.Height,
			Data:   fmt.Sprintf("This is a block data, height is %d.", s.lastBlock.Height),
		}
	}
	msg := types.ProposeMessage{
		Height:    s.currentBlock.Height,
		Validator: s.address.Id,
		Data:      s.currentBlock.Data,
		Signer:    s.address.Id,
	}
	for _, c := range s.mconn {
		if b := c.Send(0x1, msg.GetData()); !b {
			log.Println("send msg failed")
			return
		}
	}
}

func (s *State) PreVote(height int) {

}

func (s *State) PreCommit(height int) {

}

func (s *State) Commit(height int) {

}
