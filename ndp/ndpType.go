package ndp

import (
	"encoding/json"
	"github.com/CRVV/p2pFileSystem/logger"
	"github.com/CRVV/p2pFileSystem/settings"
)

var myself Peer = Peer{
	Username: settings.GetSettings().GetUsername(),
	ID:       id,
	Group:    settings.GetSettings().GetGroupName()}

type Peer struct {
	Username string
	Addr     string `json:"-"`
	ID       string
	Group    string
}

type Message struct {
	MessageType string
	Target      string
}

func (m Message) Type() string {
	return m.MessageType
}
func (m Message) Destination() string {
	var addr string
	var port string
	switch m.MessageType {
	case settings.NeighborDiscoveryProtocol:
		port = settings.NeighborDiscoveryPort
	case settings.NeighborDiscoveryProtocolEcho:
		port = settings.NeighborDiscoveryPort
	default:
		port = settings.CommunicationPort
	}

	if m.Target == settings.BroadcastAddress {
		addr = m.Target
	} else {
		addr = GetPeerAddr(m.Target)
	}
	return addr + port
}
func (m Message) Payload() []byte {
	payload, err := peer2Json(myself)
	if err != nil {
		logger.Warning(err)
		return nil
	}
	return payload
}

func peer2Json(message Peer) ([]byte, error) {
	b, err := json.Marshal(message)
	return b, err
}
func json2peer(message []byte) (Peer, error) {
	cm := Peer{}
	err := json.Unmarshal(message, &cm)
	return cm, err
}

type IUMessage struct {
	addr string
}

func (m IUMessage) Type() string {
	return settings.InvalidUsername
}
func (m IUMessage) Destination() string {
	return m.addr + settings.CommunicationPort
}
func (m IUMessage) Payload() []byte {
	return []byte(settings.InvalidUsername)
}
