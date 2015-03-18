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
	udpAddr, err := net.ResolveUDPAddr("udp", targetAddr+settings.NeighborDiscoveryPort)
	if err != nil {
		return
	}
	conn, err := net.DialUDP("udp", nil, udpAddr)
	if err != nil {
		return
	}
	defer conn.Close()

	var buff []byte
	buff = append(buff, []byte(settings.NeighborDiscoveryProtocol)...)

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
	size := make([]byte, settings.MessageHeaderSize)
	binary.LittleEndian.PutUint32(size, uint32(len(message)))
	buff = append(buff, size...)
	buff = append(buff, message...)
	_, err = conn.Write(buff)
	if err != nil {
		panic(err)
	}
}
func StartNeighborDiscoveryServer() {
	conn, err := net.ListenUDP("udp", &net.UDPAddr{
		IP:   net.IPv4zero,
		Port: 1540})
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	for {
		buff := make([]byte, settings.NeighborDiscoveryMessageBufferSize)
		_, _, err := conn.ReadFromUDP(buff)
		if err != nil {
			continue
		}
		headerSize := settings.MessageHeaderSize
		messageType := string(buff[0:headerSize])
		if messageType != settings.NeighborDiscoveryProtocol {
			continue
		}
		size := binary.LittleEndian.Uint32(buff[headerSize : headerSize*2])
		if messageType == settings.NeighborDiscoveryProtocol {
			client, err := Json2ClientMessage(buff[headerSize*2 : headerSize*2+int(size)])
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
		if ok || client.Username == settings.GetSettings().GetUsername() {
			SendMessage(client.Addr, settings.InvalidUsername, []byte(settings.InvalidUsername))
		} else {
			SendNeighborSolicitation(client.Addr)
			filesystem.OnDiscoverClient(client.Username, client.Addr)
		}
	}
}
func InitNeighborDiscovery() {
	go StartNeighborDiscoveryServer()
	timer := time.Tick(time.Second * 1)
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
