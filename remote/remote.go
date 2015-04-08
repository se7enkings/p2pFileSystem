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
		notice := <-changeNotice
		logger.Info(fmt.Sprintf("start handle a client list maintain loop %d", notice.NoticeType))
		switch notice.NoticeType {
		case ndp.ReloadPeerListNotice:
			newPeerList := ndp.GetPeerList()
			for _, name := range clients.GetNameList() {
				_, ok := newPeerList[name]
				if !ok {
					logger.Info(fmt.Sprintf("miss client %s", name))
					go onClientMissing(name)
				}
			}
			for name, _ := range newPeerList {
				if !clients.Exist(name) {
					logger.Info("found client " + name + " but do not have its file list, request it")
					onDiscoverNewClient(name)
				}
			}
		case ndp.PeerMissingNotice:
			onClientMissing(notice.PeerName)
		case ndp.NewPeerNotice:
			onDiscoverNewClient(notice.PeerName)
		}
		logger.Info(fmt.Sprintf("complete handle a client list maintain loop %d", notice.NoticeType))
	}
}
func onReceiveFilesystem(filesystemMessage []byte) {
	client, err := Json2Client(filesystemMessage)
	if err != nil {
		logger.Warning(err)
	}
	logger.Info("receive file list from " + client.Username)
	clients.AddFilesystem(client.Username, client.FileSystem)
	refreshFilesystem()
}
func onRequestedFilesystem(name string) {
	logger.Info("send filesystem to " + name)
	transfer.SendTcpMessage(&FSMessage{DestinationName: name})
}
func onClientMissing(name string) {
	logger.Info("missing client " + name)
	clients.DeleteFilesystem(name)
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
