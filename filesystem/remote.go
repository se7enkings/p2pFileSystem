package filesystem

import (
	"github.com/CRVV/p2pFileSystem/logger"
	"github.com/CRVV/p2pFileSystem/ndp"
	"github.com/CRVV/p2pFileSystem/settings"
	"github.com/CRVV/p2pFileSystem/transfer"
	"os"
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
	transfer.SendTcpMessage(&FSMessage{DestinationName: name})
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
	transfer.SendTcpMessage(&message)
}
func GetFile(hash string) error {
	FsMutex.Lock()
	// TODO: is it necessary to use a mutex on file?
	file := FileSystem[hash]
	FsMutex.Unlock()
	var blockCount int32
	lastBlockSize := int32(file.Size % settings.FileBlockSize)
	if lastBlockSize == 0 {
		blockCount = int32(file.Size / settings.FileBlockSize)
	} else {
		blockCount = int32(file.Size/settings.FileBlockSize + 1)
	}
	os.MkdirAll(settings.GetSettings().GetSharePath()+"/.temp", 0774)
	tempFile, err := os.Create(settings.GetSettings().GetSharePath() + "/.temp/" + hash)
	if err != nil {
		logger.Warning(err)
		return err
	}
	defer tempFile.Close()

	toBeCompletedBlocks := make(map[int32]bool)
	for i := int32(0); i < blockCount; i++ {
		toBeCompletedBlocks[i] = true
	}

	completeBlockNumChan := make(chan int32)
	ownerCount := len(file.Owner)
	for {
		ownerNum := 0
		runningRoutines := 0
		for blockNum, _ := range toBeCompletedBlocks {
			requestMessage := FBRMessage{
				DestinationName: file.Owner[ownerNum],
				Username:        settings.GetSettings().GetUsername(),
				FileHash:        hash,
				BlockNum:        blockNum}
			requestMessage.BlockSize = int32(settings.FileBlockSize)
			if blockNum == blockCount-1 && lastBlockSize != 0 {
				requestMessage.BlockSize = lastBlockSize
			}

			go downloadFileBlock(tempFile, &requestMessage, completeBlockNumChan)
			runningRoutines++

			ownerNum++
			if ownerNum == ownerCount {
				ownerNum = 0
			}
		}
		if runningRoutines == 0 {
			break
		}
		for i := 0; i < runningRoutines; i++ {
			blockNum := <-completeBlockNumChan
			if blockNum < 0 {
				continue
			}
			delete(toBeCompletedBlocks, blockNum)
		}
	}
	tempFile.Close()
	err = os.Rename(settings.GetSettings().GetSharePath()+"/.temp/"+hash, file.Path+"/"+file.Name)
	logger.Warning(err)

	Init()
	return nil
}
func downloadFileBlock(tempFile *os.File, requestMessage *FBRMessage, completeBlockNumChan chan int32) {
	fileDataChan := make(chan []byte)
	err := transfer.TcpConnectionForReceiveFile(requestMessage, fileDataChan)
	if err != nil {
		completeBlockNumChan <- -1
		return
	}
	fileData := <-fileDataChan
	_, err = tempFile.WriteAt(fileData, int64(requestMessage.BlockNum)*settings.FileBlockSize)
	if err != nil {
		completeBlockNumChan <- -1
		return
	}
	completeBlockNumChan <- requestMessage.BlockNum
}
func onRequestedFileBlock(requestMessage *FBRMessage) []byte {
	FslMutex.Lock()
	file, ok := FileSystemLocal[requestMessage.FileHash]
	FslMutex.Unlock()
	if !ok {
		logger.Warning("I am requested a file but I do not have it")
		return nil
	}
	f, err := os.Open(settings.GetSettings().GetSharePath()+file.Path + "/" + file.Name)
	if err != nil {
		logger.Warning(err)
		return nil
	}
	buff := make([]byte, requestMessage.BlockSize)
	_, err = f.ReadAt(buff, int64(requestMessage.BlockNum)*settings.FileBlockSize)
	if err != nil {
		logger.Warning(err)
		return nil
	}
	return buff
}
func UploadFile(path string) {

}
