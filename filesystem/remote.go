package filesystem

import (
	"encoding/json"
	"github.com/CRVV/p2pFileSystem/logger"
)

var ID string
var Clients map[string]Client

type Client struct {
	Addr       string
	Username   string
	FileSystem Filesystem
}

func OnDiscoverClient(username string, addr string) {
	Clients[username] = Client{Addr: addr, Username: username}
}
func OnReceiveFilesystem(filesystemMessage []byte) {
	client, err := Json2FileSystem(filesystemMessage)
	logger.Warning(err)
	Clients[client.Username] = client
	fsMutex.Lock()
	FileSystem = AppendFilesystem(FileSystem, client.FileSystem)
	fsMutex.Unlock()
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
