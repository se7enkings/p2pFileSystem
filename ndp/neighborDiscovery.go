package ndp

import (
	"encoding/base64"
	"encoding/binary"
	"github.com/CRVV/p2pFileSystem/logger"
	"github.com/CRVV/p2pFileSystem/settings"
	"github.com/CRVV/p2pFileSystem/transfer"
	"math/rand"
	"net"
	"reflect"
	"sync"
	"time"
)

var ClientsChangeNotice chan int = make(chan int)

var ClientList map[string]Client
var ClMutex sync.Mutex = sync.Mutex{}

var clientTemp map[string]Client // key: Username
var ctMutex sync.Mutex = sync.Mutex{}
var id string

func GetClientList() map[string]Client {
	ClMutex.Lock()
	defer ClMutex.Unlock()
	return ClientList
}
func OnExit() {
	// TODO
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
		_, remoteAddr, err := conn.ReadFromUDP(buff)
		if err != nil {
			logger.Warning(err)
			continue
		}
		headerSize := settings.MessageHeaderSize
		messageType := string(buff[0:headerSize])
		size := binary.LittleEndian.Uint32(buff[headerSize : headerSize*2])
		client, err := json2client(buff[headerSize*2 : headerSize*2+int(size)])
		if err != nil {
			logger.Warning(err)
			continue
		}
		client.Addr = remoteAddr.String()
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
func onReceiveNeighborSolicitation(client Client) {
	if client.Group == settings.GetSettings().GetGroupName() && client.ID != id {
		if client.Username == settings.GetSettings().GetUsername() {
			logger.Info("found a client which have the same username. kick it out!")
			transfer.SendTcpMessage(IUMessage{client.Addr})
		} else {
			sendNDMessage(settings.NeighborDiscoveryProtocolEcho, client.Addr)
			ClMutex.Lock()
			_, ok := ClientList[client.Username]
			if !ok {
				//				message, err := client2Json(Client{Username: settings.GetSettings().GetUsername(), Group: settings.GetSettings().GetGroupName()})
				//				if err != nil {
				//					logger.Warning(err)
				//					return
				//				}
				logger.Info("receive neighbor solicitation message from an unknown client.")
				ClientList[client.Username] = client
				ClientsChangeNotice <- 1
				//				messagePipe <- Message{Type: settings.FileSystemRequestProtocol, Destination: client.Addr, Load: message}
			}
			ClMutex.Unlock()
		}
	}
}
func onReceiveNeighborSolicitationEcho(client Client) {
	if client.Group == settings.GetSettings().GetGroupName() && client.ID != id {
		_, ok := clientTemp[client.Username]
		if ok {
			logger.Info("found a known client from " + client.Addr)
			return
		}
		logger.Info("found a new client from " + client.Addr)
		ctMutex.Lock()
		clientTemp[client.Username] = client
		ctMutex.Unlock()
	}
}
func NeighborDiscovery() {
	genID()
	for {
		doNeighborDiscovery()
		time.Sleep(time.Minute)
	}
}
func doNeighborDiscovery() {
	ctMutex.Lock()
	clientTemp = make(map[string]Client)
	ctMutex.Unlock()
	for i := 0; i < 3; i++ {
		sendNDMessage(settings.NeighborDiscoveryProtocol, settings.BroadcastAddress)
		time.Sleep(time.Second)
	}
	ctMutex.Lock()
	ClMutex.Lock()
	if !reflect.DeepEqual(ClientList, clientTemp) {
		ClientList = clientTemp
		ClientsChangeNotice <- 1
	}
	ClMutex.Unlock()
	ctMutex.Unlock()
	//	for name, _ := range filesystem.Clients {
	//		_, ok := clientTemp[name]
	//		if !ok {
	//			logger.Info("client " + name + " miss")
	//			filesystem.OnClientMissing(name)
	//		}
	//	}
	//	for name, c := range clientTemp {
	//		_, ok := filesystem.Clients[name]
	//		if !ok {
	//			message, err := client2Json(Client{Username: settings.GetSettings().GetUsername(), Group: settings.GetSettings().GetGroupName()})
	//			if err != nil {
	//				logger.Warning(err)
	//				return
	//			}
	//			logger.Info("found a client but do not have its file list, request it")
	//			messagePipe <- Message{Type: settings.FileSystemRequestProtocol, Destination: c.Addr, Load: message}
	//		}
	//	}
	//	ctMutex.Unlock()
}
func sendNDMessage(messageType string, targetAddr string) {
	message := NDMessage{myself, messageType, targetAddr}
	transfer.SendUdpPackage(message)
	logger.Info("a ndp message has been sent to " + targetAddr)
}
func genID() {
	rand.Seed(time.Now().UnixNano())
	idd := make([]byte, 16)
	for i, _ := range idd {
		idd[i] = byte(rand.Intn(256))
	}
	id = base64.StdEncoding.EncodeToString(idd)
}
