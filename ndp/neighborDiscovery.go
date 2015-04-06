package ndp

import (
	"encoding/binary"
	"fmt"
	"github.com/CRVV/p2pFileSystem/logger"
	"github.com/CRVV/p2pFileSystem/settings"
	"github.com/CRVV/p2pFileSystem/transfer"
	"net"
	"time"
)

var peersChangeNotice chan PeerListNotice

var peerList peerTable = peerTable{M: make(map[string]Peer)}
var peerListTemp peerTable = peerTable{M: make(map[string]Peer)}

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
	logger.Info(fmt.Sprintf("receive goodbye from %s", peer.Addr))
	peerList.Lock()
	delete(peerList.M, peer.Username)
	peerList.Unlock()
	peersChangeNotice <- PeerListNotice{NoticeType: PeerMissingNotice, PeerName: peer.Username}
}
func onReceiveNeighborSolicitation(peer Peer) {
	logger.Info(fmt.Sprintf("receive ndp message from %s", peer.Addr))
	if peer.Username == settings.GetSettings().GetUsername() {
		logger.Info("found a peer which have the same username. kick it out!")
		transfer.SendTcpMessage(&IUMessage{peer.Addr})
	} else {
		peerList.RLock()
		_, ok := peerList.M[peer.Username]
		peerList.RUnlock()
		if !ok {
			peerList.Lock()
			peerList.M[peer.Username] = peer
			peerList.Unlock()
		}
		sendNDMessage(settings.NeighborDiscoveryProtocolEcho, peer.Username)
		time.Sleep(time.Millisecond * 100)
		if !ok {
			peersChangeNotice <- PeerListNotice{NoticeType: NewPeerNotice, PeerName: peer.Username}
		}
	}
}
func OnReceiveNeighborSolicitationEcho(peer Peer) {
	logger.Info(fmt.Sprintf("receive ndp echo from %s", peer.Addr))
	peerList.RLock()
	_, ok := peerList.M[peer.Username]
	peerList.RUnlock()
	if !ok {
		peerList.Lock()
		peerList.M[peer.Username] = peer
		peerList.Unlock()
	}
	peerListTemp.RLock()
	_, ok = peerListTemp.M[peer.Username]
	peerListTemp.RUnlock()
	if !ok {
		peerListTemp.Lock()
		peerListTemp.M[peer.Username] = peer
		peerListTemp.Unlock()
	}
}
func NeighborDiscovery(notice chan PeerListNotice) {
	peersChangeNotice = notice
	for {
		doNeighborDiscovery()
		time.Sleep(time.Minute)
	}
}

func doNeighborDiscovery() {
	for i := 0; i < 3; i++ {
		sendNDMessage(settings.NeighborDiscoveryProtocol, settings.BroadcastAddress)
		time.Sleep(time.Second)
	}
	peerListTemp.Lock()
	peerList.Lock()
	peerList = peerListTemp
	peerListTemp.M = make(map[string]Peer)
	peerList.Unlock()
	peerListTemp.Unlock()
	peersChangeNotice <- PeerListNotice{NoticeType: ReloadPeerListNotice}
}
func sendNDMessage(messageType string, target string) {
	message := Message{messageType, target}
	transfer.SendUdpPackage(&message)
}
