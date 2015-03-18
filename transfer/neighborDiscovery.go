package transfer

import (
	"encoding/binary"
	"encoding/json"
	"github.com/CRVV/p2pFileSystem/filesystem"
	"github.com/CRVV/p2pFileSystem/settings"
	"net"
	"time"
)

func SendNeighborSolicitation(targetAddr string) {
	conn, err := net.Dial("udp", targetAddr+settings.NeighborDiscoveryPort)
	if err != nil {
		return
	}
	defer conn.Close()

	conn.Write([]byte(settings.NeighborDiscoveryProtocol))
	buff := make([]byte, settings.MessageHeaderSize)

	addr, _, err := net.SplitHostPort(conn.LocalAddr().String())
	if err != nil {
		return
	}
	message, err := ClientMessage2Json(NDMessage{
		Hello:    settings.HelloMessage,
		Username: settings.GetSettings().GetUsername(),
		ID:       filesystem.ID,
		Group:    settings.GetSettings().GetGroupName(),
		Addr:     addr})
	if err != nil {
		return
	}
	binary.LittleEndian.PutUint32(buff, uint32(len(message)))
	conn.Write(buff)
	conn.Write(message)
}
func startNeighborDiscoveryServer() {
	conn, err := net.ListenUDP("udp", &net.UDPAddr{
		IP:   net.IPv4zero,
		Port: 1540})
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	for {
		buff := make([]byte, settings.MessageHeaderSize)
		conn.Read(buff)
		messageType := string(buff)

		conn.Read(buff)
		connSize := binary.LittleEndian.Uint32(buff)

		buff = make([]byte, connSize)
		conn.Read(buff)
		if messageType == settings.NeighborDiscoveryProtocol {
			client, err := Json2ClientMessage(buff)
			if err != nil {
				continue
			}
			onReceiveNeighborSolicitation(client)
		}
	}
}
func onReceiveNeighborSolicitation(client NDMessage) {
	if client.Group == settings.GetSettings().GetGroupName() && client.ID != filesystem.ID {
		_, ok := filesystem.Clients[client.Username]
		if ok {
			SendMessage(client.Addr, settings.InvalidUsername, []byte(settings.InvalidUsername))
		} else {
			SendNeighborSolicitation(client.Addr)
			filesystem.OnDiscoverClient(client.Username, client.Addr)
		}
	}
}
func InitNeighborDiscovery() {
	go startNeighborDiscoveryServer()
	timer := time.Tick(time.Second * 5)
	for {
		SendNeighborSolicitation(settings.BroadcastAddress)
		<-timer
	}
}

type NDMessage struct {
	Hello    string
	Username string
	ID       string
	Group    string
	Addr     string
}

func ClientMessage2Json(client NDMessage) ([]byte, error) {
	b, err := json.Marshal(client)
	return b, err
}
func Json2ClientMessage(jsonClientMessage []byte) (NDMessage, error) {
	cm := NDMessage{}
	err := json.Unmarshal(jsonClientMessage, &cm)
	return cm, err
}
