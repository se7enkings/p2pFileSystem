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

var clientList map[string]NDMessage
var clMutex sync.Mutex = sync.Mutex{}

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
	message, err := NDMessage2Json(NDMessage{
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
		client, err := Json2NDMessage(buff[headerSize*2 : headerSize*2+int(size)])
		if err != nil {
			logger.Warning(err)
			continue
		}
		switch messageType {
		case settings.NeighborDiscoveryProtocol:
			logger.Info("receive ndp message from " + client.Addr)
			go onReceiveNeighborSolicitation(client)
		case settings.NeighborDiscoveryProtocolEcho:
			logger.Info("receive ndp echo message from " + client.Addr)
			go onReceiveNeighborSolicitationEcho(client)
		}
	}
}
func onReceiveNeighborSolicitation(client NDMessage) {
	if client.Group == settings.GetSettings().GetGroupName() && client.ID != id {
		switch {
		case client.Username == settings.GetSettings().GetUsername():
			logger.Info("found a client which have the same username. kick it out!")
			messagePipe <- Message{Type: settings.InvalidUsername, Load: []byte(settings.InvalidUsername), Destination: client.Addr}
		default:
			sendNeighborSolicitation(settings.NeighborDiscoveryProtocolEcho, client.Addr)
			clMutex.Lock()
			_, ok := clientList[client.Username]
			if !ok {
				message, err := NDMessage2Json(NDMessage{Username: settings.GetSettings().GetUsername(), Group: settings.GetSettings().GetGroupName()})
				if err != nil {
					logger.Warning(err)
					return
				}
				logger.Info("receive neighbor solicitation message from an unknown client, request its file list")
				clientList[client.Username] = client
				messagePipe <- Message{Type: settings.FileSystemRequestProtocol, Destination: client.Addr, Load: message}
			}
			clMutex.Unlock()
		}
	}
}
func onReceiveNeighborSolicitationEcho(client NDMessage) {
	if client.Group == settings.GetSettings().GetGroupName() && client.ID != id {
		_, ok := clients[client.Username]
		if ok {
			logger.Info("found an old client from " + client.Addr)
			return
		}
		logger.Info("found a new client from " + client.Addr)
		cMutex.Lock()
		clients[client.Username] = client
		cMutex.Unlock()
	}
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
	cMutex.Lock()
	clients = make(map[string]NDMessage)
	cMutex.Unlock()
	for i := 0; i < 3; i++ {
		sendNeighborSolicitation(settings.NeighborDiscoveryProtocol, settings.BroadcastAddress)
		time.Sleep(time.Second)
	}
	cMutex.Lock()
	clMutex.Lock()
	clientList = clients
	clMutex.Unlock()
	for name, _ := range filesystem.Clients {
		_, ok := clients[name]
		if !ok {
			logger.Info("client " + name + " miss")
			filesystem.OnClientMissing(name)
		}
	}
	for name, c := range clients {
		_, ok := filesystem.Clients[name]
		if !ok {
			message, err := NDMessage2Json(NDMessage{Username: settings.GetSettings().GetUsername(), Group: settings.GetSettings().GetGroupName()})
			if err != nil {
				logger.Warning(err)
				return
			}
			logger.Info("found a client but do not have its file list, request it")
			messagePipe <- Message{Type: settings.FileSystemRequestProtocol, Destination: c.Addr, Load: message}
		}
	}
	cMutex.Unlock()
}

type NDMessage struct {
	Username string
	Addr     string
	ID       string
	Group    string
}

func NDMessage2Json(message NDMessage) ([]byte, error) {
	b, err := json.Marshal(message)
	return b, err
}
func Json2NDMessage(jsonClientMessage []byte) (NDMessage, error) {
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
