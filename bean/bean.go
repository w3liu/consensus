package bean

type PacketMsg struct {
	ChannelID int32  `json:"channelId"`
	EOF       int32  `json:"eof"`
	Data      []byte `json:"data"`
}
