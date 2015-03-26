package ndp

import (
	"encoding/json"
	"github.com/CRVV/p2pFileSystem/logger"
	"github.com/CRVV/p2pFileSystem/settings"
)

var myself Client = Client{
	Username: settings.GetSettings().GetUsername(),
	ID:       id,
	Group:    settings.GetSettings().GetGroupName()}

type Client struct {
	Username string
	Addr     string `json:"-"`
	ID       string
	Group    string
}

type NDMessage struct {
	myselfClient Client
	messageType  string
	addr         string
}

func (m NDMessage) Type() string {
	return m.messageType
}
func (m NDMessage) Destination() string {
	return m.addr + settings.NeighborDiscoveryPort
}
func (m NDMessage) Payload() []byte {
	payload, err := client2Json(m.myselfClient)
	if err != nil {
		logger.Warning(err)
		return nil
	}
	return payload
}

func client2Json(message Client) ([]byte, error) {
	b, err := json.Marshal(message)
	return b, err
}
func json2client(jsonClientMessage []byte) (Client, error) {
	cm := Client{}
	err := json.Unmarshal(jsonClientMessage, &cm)
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
