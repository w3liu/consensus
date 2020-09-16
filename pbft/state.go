package pbft

import (
	"github.com/w3liu/consensus/config"
	"github.com/w3liu/consensus/p2p/conn"
	"time"
)

type State struct {
	peers  map[byte]*conn.MConnection
	ticker *time.Ticker // 用于定时出块
	block  Block
}

func NewState(cfg *config.Config) *State {
	return &State{
		ticker: time.NewTicker(time.Second * 10),
	}
}

func (s *State) Start() {

}
