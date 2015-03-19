package filesystem

import (
	"encoding/json"
	"github.com/CRVV/p2pFileSystem/logger"
	"github.com/CRVV/p2pFileSystem/settings"
	"sync"
)

var Clients map[string]Client
var CMutex sync.Mutex = sync.Mutex{}

var FileSystemLocal Filesystem
var FslMutex sync.Mutex = sync.Mutex{}

var MessagePipe chan Message = make(chan Message, 4)

type Client struct {
	Username   string
	FileSystem Filesystem
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
	FsMutex.Lock()
	FileSystem = AppendFilesystem(FileSystem, client.FileSystem)
	FsMutex.Unlock()
}
func OnRequestedFilesystem(name string) {
	FslMutex.Lock()
	message, err := Client2Json(Client{Username: settings.GetSettings().GetUsername(), FileSystem: FileSystemLocal})
	FslMutex.Unlock()
	if err != nil {
		logger.Warning(err)
		return
	}
	logger.Info("send filesystem to " + name)
	MessagePipe <- Message{Type: settings.FileSystemListProtocol, DestinationUsername: name, Load: message}
}
func OnClientMissing(name string) {
	logger.Info("missing client " + name)
	CMutex.Lock()
	delete(Clients, name)
	CMutex.Unlock()
	Init()
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
