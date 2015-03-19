package transfer

import (
	"encoding/binary"
	"encoding/json"
	"github.com/CRVV/p2pFileSystem/filesystem"
	"github.com/CRVV/p2pFileSystem/logger"
	"github.com/CRVV/p2pFileSystem/settings"
	"net"
)

func SendNeighborSolicitation(targetAddr string) {
	udpAddr, err := net.ResolveUDPAddr("udp", targetAddr+settings.NeighborDiscoveryPort)
	if err != nil {
		logger.Warning(err)
		return
	}

	conn, err := net.DialUDP("udp", nil, udpAddr)
	if err != nil {
		logger.Warning(err)
		return
	}
	defer conn.Close()

	var buff []byte
	buff = append(buff, []byte(settings.NeighborDiscoveryProtocol)...)

	addr, _, err := net.SplitHostPort(conn.LocalAddr().String())
	if err != nil {
		logger.Warning(err)
		return
	}
	message, err := ClientMessage2Json(NDMessage{
		Hello:    settings.HelloMessage,
		Username: settings.GetSettings().GetUsername(),
		ID:       filesystem.ID,
		Group:    settings.GetSettings().GetGroupName(),
		Addr:     addr})
	if err != nil {
		logger.Warning(err)
		return
	}
	size := make([]byte, settings.MessageHeaderSize)
	binary.LittleEndian.PutUint32(size, uint32(len(message)))
	buff = append(buff, size...)
	buff = append(buff, message...)
	_, err = conn.Write(buff)
	logger.Warning(err)
	logger.Info("a ndp message has been sent to " + targetAddr)
}
func StartNeighborDiscoveryServer() {
	conn, err := net.ListenUDP("udp", &net.UDPAddr{
		IP:   net.IPv4zero,
		Port: 1540})
	if err != nil {
		logger.Error(err)
	}
	defer conn.Close()
	for {
		buff := make([]byte, settings.NeighborDiscoveryMessageBufferSize)
		_, _, err := conn.ReadFromUDP(buff)
		if err != nil {
			logger.Warning(err)
			continue
		}
		headerSize := settings.MessageHeaderSize
		messageType := string(buff[0:headerSize])
		if messageType == settings.NeighborDiscoveryProtocol {
			size := binary.LittleEndian.Uint32(buff[headerSize : headerSize*2])
			client, err := Json2ClientMessage(buff[headerSize*2 : headerSize*2+int(size)])
			if err != nil {
				logger.Warning(err)
				continue
			}
			logger.Info("receive ndp message, " + client.String())
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
	SendNeighborSolicitation(settings.BroadcastAddress)
	//    timer := time.Tick(time.Second * 1)
	//	for {
	//		SendNeighborSolicitation(settings.BroadcastAddress)
	//		<-timer
	//	}
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
