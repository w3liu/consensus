package types

import (
	"encoding/json"
	"log"
)

type Message interface {
	GetData() []byte
}

type NodeInfoMessage struct {
	Id   string `json:"id"`
	Ip   string `json:"ip"`
	Port uint16 `json:"port"`
}

func (m *NodeInfoMessage) GetData() []byte {
	data, err := json.Marshal(m)
	if err != nil {
		log.Println(err)
		return []byte{}
	}
	return data
}

type ProposeMessage struct {
	Height    int    `json:"height"`
	Data      string `json:"data"`
	Validator string `json:"validator"`
	Signer    string `json:"signer"`
}

func (m *ProposeMessage) GetData() []byte {
	data, err := json.Marshal(m)
	if err != nil {
		log.Println(err)
		return []byte{}
	}
	return data
}

type VoteMessage struct {
	Height    int    `json:"height"`
	Validator string `json:"validator"`
}

func (m *VoteMessage) GetData() []byte {
	data, err := json.Marshal(m)
	if err != nil {
		log.Println(err)
		return []byte{}
	}
	return data
}

type PreCommitMessage struct {
	Height    int    `json:"height"`
	Validator string `json:"validator"`
}

func (m *PreCommitMessage) GetData() []byte {
	data, err := json.Marshal(m)
	if err != nil {
		log.Println(err)
		return []byte{}
	}
	return data
}

type CommitMessage struct {
	Height    int    `json:"height"`
	Validator string `json:"validator"`
}

func (m *CommitMessage) GetData() []byte {
	data, err := json.Marshal(m)
	if err != nil {
		log.Println(err)
		return []byte{}
	}
	return data
}
