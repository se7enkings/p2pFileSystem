package transfer

import (
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"github.com/CRVV/p2pFileSystem/filesystem"
	"github.com/CRVV/p2pFileSystem/logger"
	"github.com/CRVV/p2pFileSystem/settings"
	"math/rand"
	"net"
	"sync"
	"time"
)

var id string

var clients map[string]NDMessage // key: Username
var cMutex sync.Mutex = sync.Mutex{}

func sendNeighborSolicitation(messageType string, targetAddr string) {
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
	buff = append(buff, []byte(messageType)...)

	addr, _, err := net.SplitHostPort(conn.LocalAddr().String())
	if err != nil {
		logger.Warning(err)
		return
	}
	message, err := Message2Json(NDMessage{
		Username: settings.GetSettings().GetUsername(),
		ID:       id,
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
	logger.Info("start ndp server")
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
		size := binary.LittleEndian.Uint32(buff[headerSize : headerSize*2])
		client, err := Json2ClientMessage(buff[headerSize*2 : headerSize*2+int(size)])
		if err != nil {
			logger.Warning(err)
			continue
		}
		switch messageType {
		case settings.NeighborDiscoveryProtocol:
			logger.Info("receive ndp message from " + client.Addr)
			onReceiveNeighborSolicitation(client)
		case settings.NeighborDiscoveryProtocolEcho:
			logger.Info("receive ndp echo message from " + client.Addr)
			onReceiveNeighborSolicitationEcho(client)
		}
	}
}
func onReceiveNeighborSolicitation(client NDMessage) {
	if client.Group == settings.GetSettings().GetGroupName() && client.ID != id {
		cMutex.Lock()
		_, ok := clients[client.Username]
		cMutex.Unlock()
		switch {
		case client.Username == settings.GetSettings().GetUsername():
			logger.Info("found a client which have the same username. kick it out!")
			messagePipe <- Message{Type: settings.InvalidUsername, Load: []byte(settings.InvalidUsername), Destination: client.Addr}
		case ok:
			logger.Info("found a known client")
			sendNeighborSolicitation(settings.NeighborDiscoveryProtocolEcho, client.Addr)
		default:
			logger.Info("found an unknown client " + client.Username + " from " + client.Addr)
			sendNeighborSolicitation(settings.NeighborDiscoveryProtocolEcho, client.Addr)
		}
	}
}
func onReceiveNeighborSolicitationEcho(client NDMessage) {
	cMutex.Lock()
	clients[client.Username] = client
	cMutex.Unlock()
}
func InitNeighborDiscovery() {
	genID()
	StartNeighborDiscoveryServer()
}
func NeighborDiscovery() {
	for {
		DoNeighborDiscovery()
		time.Sleep(time.Minute)
	}
}
func DoNeighborDiscovery() {
	clients = make(map[string]NDMessage)
	for i := 0; i < 3; i++ {
		sendNeighborSolicitation(settings.NeighborDiscoveryProtocol, settings.BroadcastAddress)
		time.Sleep(time.Second * 2)
	}
	for name, _ := range filesystem.Clients {
		_, ok := clients[name]
		if !ok {
			filesystem.OnClientMissing(name)
		}
	}
	for name, c := range clients {
		_, ok := filesystem.Clients[name]
		if !ok {
			messagePipe <- Message{Type: settings.FileSystemRequestProtocol, Destination: c.Addr, Load: []byte(settings.GetSettings().GetUsername())}
		}
	}
}

type NDMessage struct {
	Username string
	Addr     string
	ID       string
	Group    string
}

func Message2Json(message interface{}) ([]byte, error) {
	b, err := json.Marshal(message)
	return b, err
}
func Json2ClientMessage(jsonClientMessage []byte) (NDMessage, error) {
	cm := NDMessage{}
	err := json.Unmarshal(jsonClientMessage, &cm)
	return cm, err
}
func genID() {
	rand.Seed(time.Now().UnixNano())
	idd := make([]byte, 16)
	for i, _ := range idd {
		idd[i] = byte(rand.Intn(256))
	}
	id = base64.StdEncoding.EncodeToString(idd)
}
