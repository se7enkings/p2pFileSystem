package filesystem

import (
	"encoding/json"
	"github.com/CRVV/p2pFileSystem/logger"
	"github.com/CRVV/p2pFileSystem/ndp"
	"github.com/CRVV/p2pFileSystem/settings"
)

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
	FslMutex.Lock()
	myself := Client{Username: settings.GetSettings().GetUsername(), FileSystem: FileSystemLocal}
	message, err := Struct2Json(&myself)
	FslMutex.Unlock()
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
