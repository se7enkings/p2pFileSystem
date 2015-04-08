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
	udpAddr, err := net.ResolveUDPAddr("udp4", settings.NeighborDiscoveryPort)
	logger.Error(err)
	conn, err := net.ListenUDP("udp4", udpAddr)

	logger.Error(err)
	defer conn.Close()
	for {
		buff := make([]byte, settings.MessageBufferSize)
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
	peerList.Delete(peer.Username)
	peersChangeNotice <- PeerListNotice{NoticeType: PeerMissingNotice, PeerName: peer.Username}
}
func onReceiveNeighborSolicitation(peer Peer) {
	logger.Info(fmt.Sprintf("receive ndp message from %s at %s", peer.Username, peer.Addr))
	if peer.Username == settings.GetSettings().GetUsername() {
		logger.Info("found a peer which have the same username. kick it out!")
		transfer.SendTcpMessage(&IUMessage{peer.Addr})
	} else {
		ok := peerList.Exist(peer.Username)
		if !ok {
			peerList.Add(peer)
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
	ok := peerList.Exist(peer.Username)
	if !ok {
		peerList.Add(peer)
	}
	ok = peerListTemp.Exist(peer.Username)
	if !ok {
		peerListTemp.Add(peer)
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
	for name, peer := range peerList.GetMap() {
		message := Message{settings.NeighborDiscoveryProtocol, name}
		err := transfer.SendTcpMessage(&message)
		if err == nil {
			peerListTemp.Add(peer)
		}
	}
	mapTemp := peerListTemp.GetMap()
	peerList.ReplaceByNewMap(mapTemp)
	peerListTemp.ReplaceByNewMap(make(map[string]Peer))
	peersChangeNotice <- PeerListNotice{NoticeType: ReloadPeerListNotice}
}
func sendNDMessage(messageType string, target string) {
	message := Message{messageType, target}
	transfer.SendUdpPackage(&message)
}
