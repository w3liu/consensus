package p2p

type NodeInfo interface {
	ID() ID
	nodeInfoAddress
	nodeInfoTransport
}

type nodeInfoAddress interface {
	NetAddress() (*NetAddress, error)
}

type nodeInfoTransport interface {
	Validate() error
	CompatibleWith(other NodeInfo) error
}

type DefaultNodeInfo struct {
	DefaultNodeID ID     `json:"id"`
	ListenAddr    string `json:"listen_addr"`
	Network       string `json:"network"`
}

func (s DefaultNodeInfo) ID() ID {
	return s.DefaultNodeID
}

func (s DefaultNodeInfo) Validate() error {
	return nil
}

func (s DefaultNodeInfo) CompatibleWith(otherInfo NodeInfo) error {
	return nil
}

func (s DefaultNodeInfo) NetAddress() (*NetAddress, error) {
	idAddr := IDAddressString(s.ID(), s.ListenAddr)
	return NewNetAddressString(idAddr)
}

func NewDefaultNodeInfo(id, addr, network string) DefaultNodeInfo {
	dni := DefaultNodeInfo{
		DefaultNodeID: ID(id),
		ListenAddr:    addr,
		Network:       network,
	}
	return dni
}
