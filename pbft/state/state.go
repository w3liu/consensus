package state

import (
	"encoding/json"
	"fmt"
	"github.com/w3liu/consensus/libs/gobio"
	"github.com/w3liu/consensus/log"
	"github.com/w3liu/consensus/p2p/conn"
	"github.com/w3liu/consensus/pbft/config"
	"github.com/w3liu/consensus/types"
	"go.uber.org/zap"
	"net"
	"sync"
	"time"
)

type State struct {
	address        *types.Address
	proposer       *types.Address
	validators     []*types.Address
	mconn          map[string]*conn.MConnection
	dialTicker     *time.Ticker
	blockTicker    *time.Ticker // 用于定时出块
	lastBlock      Block
	currentBlock   Block
	step           int // 0 pending, 1 propose, 2 vote, 3 precommit, 4 commit
	lock           sync.Mutex
	voteCache      map[string]*types.VoteMessage
	preCommitCache map[string]*types.PreCommitMessage
	commitCache    map[string]*types.CommitMessage
}

func NewState(cfg *config.Config) *State {
	address, err := types.NewAddress(cfg.Peer.Address)
	if err != nil {
		panic(err)
	}
	if len(cfg.Peer.Seeds) == 0 {
		panic("len(cfg.Peer.Seeds) == 0")
	}
	proposer, err := types.NewAddress(cfg.Peer.Seeds[0])
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
		proposer:       proposer,
		address:        address,
		validators:     seeds,
		mconn:          make(map[string]*conn.MConnection),
		dialTicker:     time.NewTicker(time.Second * 10),
		blockTicker:    time.NewTicker(time.Second * 10),
		voteCache:      make(map[string]*types.VoteMessage),
		preCommitCache: make(map[string]*types.PreCommitMessage),
	}
}

func (s *State) Start() {
	go func() {
		s.dialPeers()
	}()
	go func() {
		s.accept()
	}()
	for {
		select {
		case <-s.blockTicker.C:
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
	for _, addr := range s.validators {
		if _, ok := s.mconn[addr.Id]; !ok {
			c, err := s.dial(addr)
			if err != nil {
				log.Info("s.dial(addr) error", zap.Error(err))
				continue
			}
			err = c.OnStart()
			if err != nil {
				log.Info("c.OnStart()", zap.Error(err))
			}
		}
	}
	time.Sleep(time.Second * 10)
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
		log.Error("handshake error", zap.Error(err))
		_ = c.Close()
		return nil, err
	}
	if _, ok := s.mconn[peer.Id]; ok {
		log.Warn("connection is existed", zap.Any("peer", peer))
		return nil, err
	}

	onError := func(r interface{}) {
		// 移除节点
		if _, ok := s.mconn[peer.Id]; ok {
			delete(s.mconn, peer.Id)
		}
		log.Error("onError", zap.Any("r", r))
	}
	chDescs := []*conn.ChannelDescriptor{{ID: 0x01, Priority: 1, SendQueueCapacity: 1}}
	server := conn.NewMConnection(c, chDescs, s.ReceiveMsg, onError)
	s.mconn[peer.Id] = server
	return server, nil
}

func (s *State) ReceiveMsg(chID byte, msgBytes []byte) {
	log.Info("ReceiveMsg:", zap.Any("chID", chID), zap.String("msg", string(msgBytes)))
	msgInfo := &types.MessageInfo{}
	err := json.Unmarshal(msgBytes, msgInfo)
	if err != nil {
		log.Error("json.Unmarshal", zap.Error(err))
		return
	}
	switch msgInfo.MsgType {
	case types.Propose:
		proposeMsg := &types.ProposeMessage{}
		err := json.Unmarshal([]byte(msgInfo.MsgContent), proposeMsg)
		if err != nil {
			log.Error("json.Unmarshal", zap.Error(err))
			return
		}
		s.ReceivePropose(proposeMsg)
	case types.Vote:
		voteMsg := &types.VoteMessage{}
		err := json.Unmarshal([]byte(msgInfo.MsgContent), voteMsg)
		if err != nil {
			log.Error("json.Unmarshal", zap.Error(err))
			return
		}
		s.ReceiveVote(voteMsg)
	case types.PreCommit:
		preCommitMsg := &types.PreCommitMessage{}
		err := json.Unmarshal([]byte(msgInfo.MsgContent), preCommitMsg)
		if err != nil {
			log.Error("json.Unmarshal", zap.Error(err))
			return
		}
		s.ReceivePreCommit(preCommitMsg)

	}
}

func (s *State) Propose() {
	if s.step == types.Pending && s.isProposer() {
		currHeight := s.lastBlock.Height + 1
		s.currentBlock = Block{
			Height: currHeight,
			Data:   fmt.Sprintf("This is a block data, height is %d.", currHeight),
		}
	} else {
		return
	}

	if !s.checkMajor23(len(s.mconn) + 1) {
		log.Warn("connected peer number is lower than 2/3")
		return
	}
	msg := &types.ProposeMessage{
		Height:    s.currentBlock.Height,
		Validator: s.address.Id,
		Data:      s.currentBlock.Data,
		Signer:    s.address.Id,
	}
	for _, c := range s.mconn {
		if b := c.Send(0x1, types.NewMessageInfo(msg).GetData()); !b {
			log.Error("send msg failed")
			return
		}
	}
	s.step = types.Propose
}

func (s *State) Vote() {
	if s.step == types.Pending {
		msg := &types.VoteMessage{
			Height:    s.currentBlock.Height,
			Validator: s.address.Id,
		}
		for _, c := range s.mconn {
			if b := c.Send(0x1, types.NewMessageInfo(msg).GetData()); !b {
				log.Error("send msg failed")
				return
			}
		}
		s.step = types.Vote
	}
}

func (s *State) PreCommit() {
	if s.step == types.Vote {
		msg := &types.PreCommitMessage{
			Height:    s.currentBlock.Height,
			Validator: s.address.Id,
		}
		for _, c := range s.mconn {
			if b := c.Send(0x1, types.NewMessageInfo(msg).GetData()); !b {
				log.Error("send msg failed")
				return
			}
		}
		s.step = types.PreCommit
	}

}

func (s *State) Commit() {
	if s.step == types.PreCommit {
		s.lastBlock = s.currentBlock
		s.currentBlock = Block{}
		s.step = types.Pending
		log.Info("block commit", zap.Any("height", s.lastBlock.Height), zap.Any("data", s.lastBlock.Data))
	}
}

func (s *State) ReceivePropose(o *types.ProposeMessage) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.currentBlock = Block{
		Height: o.Height,
		Data:   o.Data,
	}
	if s.checkBlock() {
		s.Vote()
	} else {
		s.currentBlock = s.lastBlock
	}
}

func (s *State) ReceiveVote(o *types.VoteMessage) {
	s.lock.Lock()
	defer s.lock.Unlock()
	if _, ok := s.voteCache[o.Validator]; !ok {
		s.voteCache[o.Validator] = o
	}
	var voteCnt int
	if s.isProposer() {
		voteCnt = len(s.voteCache) + 1
	} else {
		voteCnt = len(s.voteCache) + 2
	}
	if s.checkMajor23(voteCnt) {
		s.PreCommit()
	}
}

func (s *State) ReceivePreCommit(o *types.PreCommitMessage) {
	s.lock.Lock()
	defer s.lock.Unlock()
	if _, ok := s.preCommitCache[o.Validator]; !ok {
		s.preCommitCache[o.Validator] = o
	}
	var pCnt = len(s.preCommitCache) + 1
	if s.checkMajor23(pCnt) {
		s.Commit()
	}
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

func (s *State) isProposer() bool {
	return s.proposer.Id == s.address.Id
}

func (s *State) checkBlock() bool {
	// TODO:检查签名
	return true
}

func (s *State) checkMajor23(cnt int) bool {
	n := len(s.validators) + 1
	return (cnt * 1000 / n) > 2*1000/3
}
