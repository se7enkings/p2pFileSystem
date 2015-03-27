package ndp

import (
	"encoding/base64"
	"encoding/binary"
	"github.com/CRVV/p2pFileSystem/logger"
	"github.com/CRVV/p2pFileSystem/settings"
	"github.com/CRVV/p2pFileSystem/transfer"
	"math/rand"
	"net"
	"sync"
	"time"
)

var peersChangeNotice chan string

var peerList map[string]Peer = make(map[string]Peer)
var plMutex sync.Mutex = sync.Mutex{}

var peerListTemp map[string]Peer = make(map[string]Peer) // key: Username
var pltMutex sync.Mutex = sync.Mutex{}
var id string

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
		peer, err := json2peer(buff[headerSize*2 : headerSize*2+int(size)])
		if err != nil {
			logger.Warning(err)
			continue
		}
		peer.Addr, _, _ = net.SplitHostPort(remoteAddr.String())
		switch messageType {
		case settings.NeighborDiscoveryProtocol:
			go onReceiveNeighborSolicitation(peer)
		case settings.NeighborDiscoveryProtocolEcho:
			go OnReceiveNeighborSolicitationEcho(peer)
		}
	}
}
func onReceiveNeighborSolicitation(peer Peer) {
	if peer.Group == settings.GetSettings().GetGroupName() && peer.ID != id {
		if peer.Username == settings.GetSettings().GetUsername() {
			logger.Info("found a peer which have the same username. kick it out!")
			transfer.SendTcpMessage(IUMessage{peer.Addr})
		} else {
			plMutex.Lock()
			_, ok := peerList[peer.Username]
			if !ok {
				logger.Info("receive neighbor solicitation message from an unknown client " + peer.Username)
				peerList[peer.Username] = peer
			}
			plMutex.Unlock()
			sendNDMessage(settings.NeighborDiscoveryProtocolEcho, peer.Username)
			time.Sleep(time.Millisecond * 100)
			if !ok {
				peersChangeNotice <- peer.Username
			}
		}
	}
}
func OnReceiveNeighborSolicitationEcho(peer Peer) {
	if peer.Group == settings.GetSettings().GetGroupName() && peer.ID != id {
		pltMutex.Lock()
		_, ok := peerListTemp[peer.Username]
		defer pltMutex.Unlock()
		if !ok {
			logger.Info("found a new peer from " + peer.Addr)
			peerListTemp[peer.Username] = peer
		}
	}
}
func NeighborDiscovery(notice chan string) {
	peersChangeNotice = notice
	genID()
	for {
		doNeighborDiscovery()
		time.Sleep(time.Minute)
	}
}

const ReloadPeerList string = "reloaded peer list"

func doNeighborDiscovery() {
	pltMutex.Lock()
	peerListTemp = make(map[string]Peer)
	pltMutex.Unlock()
	for i := 0; i < 3; i++ {
		sendNDMessage(settings.NeighborDiscoveryProtocol, settings.BroadcastAddress)
		time.Sleep(time.Second)
	}
	pltMutex.Lock()
	plMutex.Lock()
	peerList = peerListTemp
	peersChangeNotice <- ReloadPeerList
	plMutex.Unlock()
	pltMutex.Unlock()
}
func sendNDMessage(messageType string, target string) {
	message := Message{messageType, target}
	transfer.SendUdpPackage(message)
}
func genID() {
	//	logger.Info("generate a new ID")
	rand.Seed(time.Now().UnixNano())
	idd := make([]byte, 16)
	for i, _ := range idd {
		idd[i] = byte(rand.Intn(256))
	}
	id = base64.StdEncoding.EncodeToString(idd)
	myself.ID = id
}
