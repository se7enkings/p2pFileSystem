package filesystem

import (
	"encoding/json"
)

var Clients map[string]Client

type Client struct {
	Addr       string
	Username   string
	FileSystem Filesystem
}

func OnDiscoverClient(username string, addr string) {

}
func OnReceiveFilesystem(filesystemMessage []byte) {

}
func FileSystem2Json(fileSystem Filesystem) ([]byte, error) {
	b, err := json.Marshal(fileSystem)
	return b, err
}
func Json2FileSystem(jsonFileListMessage []byte) (Filesystem, error) {
	fs := make(Filesystem)
	err := json.Unmarshal(jsonFileListMessage, &fs)
	return fs, err
}
