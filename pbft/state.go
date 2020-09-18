package pbft

import (
	"fmt"
	"github.com/w3liu/consensus/config"
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
			onReceive := func(chID byte, msgBytes []byte) {
			}
			onError := func(r interface{}) {
			}
			chDescs := []*conn.ChannelDescriptor{{ID: 0x01, Priority: 1, SendQueueCapacity: 1}}
			server := conn.NewMConnection(c, chDescs, onReceive, onError)

			err := server.OnStart()
			if err != nil {

				return
			}

		}(c)
	}
}

func (s *State) CreateConn() {

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
