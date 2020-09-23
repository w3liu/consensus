package types

import (
	"encoding/json"
	"github.com/w3liu/consensus/log"
	"go.uber.org/zap"
)

const (
	Pending = iota
	Propose
	Vote
	PreCommit
	Commit
)

type Message interface {
	GetData() []byte
}

type Serializable interface {
	Serialize() []byte
}

type NodeInfoMessage struct {
	Id   string `json:"id"`
	Ip   string `json:"ip"`
	Port uint16 `json:"port"`
}

func (m *NodeInfoMessage) GetData() []byte {
	data, err := json.Marshal(m)
	if err != nil {
		log.Error("json.Marshal", zap.Error(err))
		return []byte{}
	}
	return data
}

type MessageInfo struct {
	MsgType    int    `json:"msgType"`
	MsgContent string `json:"msgContent"`
}

func (m *MessageInfo) GetData() []byte {
	data, err := json.Marshal(m)
	if err != nil {
		log.Error("json.Marshal", zap.Error(err))
		return []byte{}
	}
	return data
}

func NewMessageInfo(o Serializable) *MessageInfo {
	msgInfo := &MessageInfo{
		MsgContent: string(o.Serialize()),
	}
	switch interface{}(o).(type) {
	case *ProposeMessage:
		msgInfo.MsgType = Propose
	case *VoteMessage:
		msgInfo.MsgType = Vote
	case *PreCommitMessage:
		msgInfo.MsgType = PreCommit
	case *CommitMessage:
		msgInfo.MsgType = Commit
	default:

	}
	return msgInfo
}

type ProposeMessage struct {
	Height    int    `json:"height"`
	Data      string `json:"data"`
	Validator string `json:"validator"`
	Signer    string `json:"signer"`
}

func (m *ProposeMessage) Serialize() []byte {
	data, err := json.Marshal(m)
	if err != nil {
		log.Error("json.Marshal", zap.Error(err))
		return []byte{}
	}
	return data
}

type VoteMessage struct {
	Height    int    `json:"height"`
	Validator string `json:"validator"`
}

func (m *VoteMessage) Serialize() []byte {
	data, err := json.Marshal(m)
	if err != nil {
		log.Error("json.Marshal", zap.Error(err))
		return []byte{}
	}
	return data
}

type PreCommitMessage struct {
	Height    int    `json:"height"`
	Validator string `json:"validator"`
}

func (m *PreCommitMessage) Serialize() []byte {
	data, err := json.Marshal(m)
	if err != nil {
		log.Error("json.Marshal", zap.Error(err))
		return []byte{}
	}
	return data
}

type CommitMessage struct {
	Height    int    `json:"height"`
	Validator string `json:"validator"`
}

func (m *CommitMessage) Serialize() []byte {
	data, err := json.Marshal(m)
	if err != nil {
		log.Error("json.Marshal", zap.Error(err))
		return []byte{}
	}
	return data
}
