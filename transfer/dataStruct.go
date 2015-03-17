package transfer

import (
	"encoding/json"
	"github.com/CRVV/p2pFileSystem/filesystem"
)

func FileSystem2Json(fileSystem filesystem.Filesystem) ([]byte, error) {
	b, err := json.Marshal(fileSystem)
	return b, err
}

func Json2FileSystem(jsonFileListMessage []byte) (filesystem.Filesystem, error) {
	fs := make(filesystem.Filesystem)
	err := json.Unmarshal(jsonFileListMessage, &fs)
	return fs, err
}

func ClientMessage2Json(client Client) ([]byte, error) {
    b, err := json.Marshal(client)
    return b, err
}
func Json2ClientMessage(jsonClientMessage []byte) (Client, error) {
    cm := Client{}
    err := json.Unmarshal(jsonClientMessage, &cm)
    return cm, err
}