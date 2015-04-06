package remote

import (
	"encoding/json"
	"github.com/CRVV/p2pFileSystem/filesystem"
	"github.com/CRVV/p2pFileSystem/logger"
	"github.com/CRVV/p2pFileSystem/ndp"
	"github.com/CRVV/p2pFileSystem/settings"
	"sync"
)

type Client struct {
	Username   string
	FileSystem filesystem.Filesystem
}
type ClientList struct {
	M map[string]filesystem.Filesystem
	sync.RWMutex
}

type FSMessage struct {
	DestinationName string
}

func (m *FSMessage) Type() string {
	return settings.FileSystemListProtocol
}
func (m *FSMessage) Destination() string {
	return ndp.GetPeerAddr(m.DestinationName)
}
func (m *FSMessage) Payload() []byte {
	fsl := filesystem.GetLocalFilesystemForSend()
	fsl.RLock()
	defer fsl.RUnlock()
	myself := Client{Username: settings.GetSettings().GetUsername(), FileSystem: *fsl}
	message, err := Struct2Json(&myself)
	if err != nil {
		logger.Warning(err)
		return nil
	}
	return message
}
func Struct2Json(fileSystem interface{}) ([]byte, error) {
	b, err := json.Marshal(fileSystem)
	return b, err
}
func Json2Client(jsonFileListMessage []byte) (Client, error) {
	fs := Client{}
	err := json.Unmarshal(jsonFileListMessage, &fs)
	return fs, err
}

type FBRMessage struct {
	DestinationName string `json:"-"`
	Username        string
	FileHash        string
	BlockNum        int32
	BlockSize       int32
}

func (m *FBRMessage) Type() string {
	return settings.FileBlockRequestProtocol
}
func (m *FBRMessage) Destination() string {
	return ndp.GetPeerAddr(m.DestinationName)
}
func (m *FBRMessage) Payload() []byte {
	message, err := Struct2Json(m)
	if err != nil {
		logger.Warning(err)
		return nil
	}
	return message
}
func Json2FBRMessage(jsonFileListMessage []byte) (FBRMessage, error) {
	fs := FBRMessage{}
	err := json.Unmarshal(jsonFileListMessage, &fs)
	return fs, err
}
