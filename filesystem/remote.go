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
	client, err := Json2FileSystem(filesystemMessage)
	if err != nil {
		logger.Warning(err)
	}
	CMutex.Lock()
	Clients[client.Username] = client
	CMutex.Unlock()
	FsMutex.Lock()
	FileSystem = AppendFilesystem(FileSystem, client.FileSystem)
	FsMutex.Unlock()
}
func OnRequestedFilesystem(name string) {
	FslMutex.Lock()
	message, err := FileSystem2Json(Client{Username: settings.GetSettings().GetUsername(), FileSystem: FileSystemLocal})
	FslMutex.Unlock()
	if err != nil {
		logger.Warning(err)
		return
	}
	MessagePipe <- Message{Type: settings.FileSystemListProtocol, DestinationUsername: name, Load: message}
}
func OnClientMissing(name string) {
	CMutex.Lock()
	delete(Clients, name)
	CMutex.Unlock()
	Init()
}

func FileSystem2Json(fileSystem Client) ([]byte, error) {
	b, err := json.Marshal(fileSystem)
	return b, err
}
func Json2FileSystem(jsonFileListMessage []byte) (Client, error) {
	fs := Client{}
	err := json.Unmarshal(jsonFileListMessage, &fs)
	return fs, err
}
