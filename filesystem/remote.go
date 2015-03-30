package filesystem

import (
	"encoding/json"
	"github.com/CRVV/p2pFileSystem/logger"
	"github.com/CRVV/p2pFileSystem/ndp"
	"github.com/CRVV/p2pFileSystem/settings"
	"github.com/CRVV/p2pFileSystem/transfer"
	"sync"
)

var Clients map[string]Client
var CMutex sync.Mutex = sync.Mutex{}

var FileSystemLocal Filesystem
var FslMutex sync.Mutex = sync.Mutex{}

type Client struct {
	Username   string
	FileSystem Filesystem
}

func MaintainClientList() {
	changeNotice := make(chan ndp.PeerListNotice)
	go ndp.NeighborDiscovery(changeNotice)
	for {
		switch notice := <-changeNotice; notice.NoticeType {
		case ndp.ReloadPeerListNotice:
			newPeerList := ndp.GetPeerList()
			for name, _ := range Clients {
				_, ok := newPeerList[name]
				if !ok {
					logger.Info("client " + name + " miss")
					onClientMissing(name)
				}
			}
			for name, _ := range newPeerList {
				_, ok := Clients[name]
				if !ok {
					logger.Info("found client " + name + " but do not have its file list, request it")
					onDiscoverNewClient(name)
				}
			}
		case ndp.PeerMissingNotice:
			onClientMissing(notice.PeerName)
		case ndp.NewPeerNotice:
			onDiscoverNewClient(notice.PeerName)
		}
	}
}
func OnReceiveFilesystem(filesystemMessage []byte) {
	client, err := Json2Client(filesystemMessage)
	if err != nil {
		logger.Warning(err)
	}
	logger.Info("receive file list from " + client.Username)
	CMutex.Lock()
	Clients[client.Username] = client
	CMutex.Unlock()
	Init()
}
func OnRequestedFilesystem(name string) {
	logger.Info("send filesystem to " + name)
	transfer.SendTcpMessage(FSMessage{DestinationName: name})
}
func onClientMissing(name string) {
	logger.Info("missing client " + name)
	CMutex.Lock()
	delete(Clients, name)
	CMutex.Unlock()
	Init()
}
func onDiscoverNewClient(name string) {
	message := ndp.Message{MessageType: settings.FileSystemRequestProtocol, Target: name}
	transfer.SendTcpMessage(message)
}
func GetFile(hash string) {
	FsMutex.Lock()
	fileOwnerList := FileSystem[hash].Owner
	FsMutex.Unlock()

}
func UploadFile(path string) {

}
func Client2Json(fileSystem Client) ([]byte, error) {
	b, err := json.Marshal(fileSystem)
	return b, err
}
func Json2Client(jsonFileListMessage []byte) (Client, error) {
	fs := Client{}
	err := json.Unmarshal(jsonFileListMessage, &fs)
	return fs, err
}
