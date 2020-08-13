package bean

type PacketMsg struct {
	ChannelID int32
	EOF       int32
	Data      []byte
}

func (m *PacketMsg) GetData() []byte {
	if m != nil {
		return m.Data
	}
	return nil
}
