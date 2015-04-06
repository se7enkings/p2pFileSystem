package ndp

import (
	"encoding/binary"
	"github.com/CRVV/p2pFileSystem/logger"
	"github.com/CRVV/p2pFileSystem/settings"
	"github.com/CRVV/p2pFileSystem/transfer"
	"net"
	"time"
)

var peersChangeNotice chan PeerListNotice

var peerList peerTable = peerTable{m: make(map[string]Peer)}
var peerListTemp peerTable = peerTable{m: make(map[string]Peer)}

func OnExit() {
	sendNDMessage(settings.GoodByeProtocol, settings.BroadcastAddress)
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
		if peer.Group != settings.GetSettings().GetGroupName() || peer.ID == myself.ID {
			continue
		}
		peer.Addr, _, _ = net.SplitHostPort(remoteAddr.String())
		switch messageType {
		case settings.NeighborDiscoveryProtocol:
			go onReceiveNeighborSolicitation(peer)
		case settings.NeighborDiscoveryProtocolEcho:
			go OnReceiveNeighborSolicitationEcho(peer)
		case settings.GoodByeProtocol:
			go onMissingPeer(peer)
		}
	}
}
func onMissingPeer(peer Peer) {
	logger.Info("receive goodbye from " + peer.Username)
	plMutex.Lock()
	delete(peerList, peer.Username)
	plMutex.Unlock()
	peersChangeNotice <- PeerListNotice{NoticeType: PeerMissingNotice, PeerName: peer.Username}
}
func onReceiveNeighborSolicitation(peer Peer) {
	if peer.Username == settings.GetSettings().GetUsername() {
		logger.Info("found a peer which have the same username. kick it out!")
		transfer.SendTcpMessage(&IUMessage{peer.Addr})
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
			peersChangeNotice <- PeerListNotice{NoticeType: NewPeerNotice, PeerName: peer.Username}
		}
	}
}
func OnReceiveNeighborSolicitationEcho(peer Peer) {
	plMutex.Lock()
	_, ok := peerList[peer.Username]
	if !ok {
		peerList[peer.Username] = peer
		logger.Info("found a new peer from " + peer.Addr)
	}
	plMutex.Unlock()
	pltMutex.Lock()
	_, ok = peerListTemp[peer.Username]
	if !ok {
		peerListTemp[peer.Username] = peer
	}
	pltMutex.Unlock()
}
func NeighborDiscovery(notice chan PeerListNotice) {
	peersChangeNotice = notice
	for {
		doNeighborDiscovery()
		time.Sleep(time.Minute)
	}
}

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
	peersChangeNotice <- PeerListNotice{NoticeType: ReloadPeerListNotice}
	plMutex.Unlock()
	pltMutex.Unlock()
}
func sendNDMessage(messageType string, target string) {
	message := Message{messageType, target}
	transfer.SendUdpPackage(&message)
}
