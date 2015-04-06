package ndp

import (
	"github.com/CRVV/p2pFileSystem/logger"
	"github.com/CRVV/p2pFileSystem/settings"
)

type Message struct {
	MessageType string
	Target      string
}

func (m *Message) Type() string {
	return m.MessageType
}
func (m *Message) Destination() string {
	var addr string

	if m.Target == settings.BroadcastAddress {
		addr = m.Target
	} else {
		var err error
		addr, err = GetPeerAddr(m.Target)
		logger.Error(err)
	}
	return addr
}
func (m *Message) Payload() []byte {
	payload, err := peer2Json(myself)
	if err != nil {
		logger.Warning(err)
		return nil
	}
	return payload
}

type IUMessage struct {
	addr string
}

func (m *IUMessage) Type() string {
	return settings.InvalidUsername
}
func (m *IUMessage) Destination() string {
	return m.addr
}
func (m *IUMessage) Payload() []byte {
	return []byte(settings.InvalidUsername)
}
