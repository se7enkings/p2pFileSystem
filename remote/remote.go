package remote

import (
	"fmt"
	"github.com/CRVV/p2pFileSystem/filesystem"
	"github.com/CRVV/p2pFileSystem/logger"
	"github.com/CRVV/p2pFileSystem/ndp"
	"github.com/CRVV/p2pFileSystem/settings"
	"github.com/CRVV/p2pFileSystem/transfer"
)

var clients ClientList = ClientList{M: make(map[string]filesystem.Filesystem)}

func MaintainClientList() {
	changeNotice := make(chan ndp.PeerListNotice)
	go ndp.NeighborDiscovery(changeNotice)
	for {
		switch notice := <-changeNotice; notice.NoticeType {
		case ndp.ReloadPeerListNotice:
			newPeerList := ndp.GetPeerList()
			newPeerList.RLock()
			clients.RLock()
			for name, _ := range clients.M {
				_, ok := newPeerList.M[name]
				if !ok {
					logger.Info(fmt.Sprintf("miss client %s", name))
					onClientMissing(name)
				}
			}
			for name, _ := range newPeerList.M {
				_, ok := clients.M[name]
				if !ok {
					logger.Info("found client " + name + " but do not have its file list, request it")
					onDiscoverNewClient(name)
				}
			}
			clients.RUnlock()
			newPeerList.RUnlock()
		case ndp.PeerMissingNotice:
			onClientMissing(notice.PeerName)
		case ndp.NewPeerNotice:
			onDiscoverNewClient(notice.PeerName)
		}
	}
}
func onReceiveFilesystem(filesystemMessage []byte) {
	client, err := Json2Client(filesystemMessage)
	if err != nil {
		logger.Warning(err)
	}
	logger.Info("receive file list from " + client.Username)
	clients.Lock()
	clients.M[client.Username] = client.FileSystem
	clients.Unlock()
	refreshFilesystem()
}
func onRequestedFilesystem(name string) {
	logger.Info("send filesystem to " + name)
	transfer.SendTcpMessage(&FSMessage{DestinationName: name})
}
func onClientMissing(name string) {
	logger.Info("missing client " + name)
	clients.Lock()
	delete(clients.M, name)
	clients.Unlock()
	refreshFilesystem()
}
func onDiscoverNewClient(name string) {
	message := ndp.Message{MessageType: settings.FileListRequestProtocol, Target: name}
	transfer.SendTcpMessage(&message)
}
func refreshFilesystem() {
	clients.RLock()
	defer clients.RUnlock()
	filesystem.RefreshRemoteFile(clients.M)
}
