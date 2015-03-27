package filesystem

import (
	"github.com/CRVV/p2pFileSystem/logger"
	"github.com/CRVV/p2pFileSystem/ndp"
	"github.com/CRVV/p2pFileSystem/settings"
)

type FSMessage struct {
	DestinationName string
}

func (m FSMessage) Type() string {
	return settings.FileSystemListProtocol
}
func (m FSMessage) Destination() string {
	return ndp.GetPeerAddr(m.DestinationName)
}
func (m FSMessage) Payload() []byte {
	FslMutex.Lock()
	message, err := Client2Json(Client{Username: settings.GetSettings().GetUsername(), FileSystem: FileSystemLocal})
	FslMutex.Unlock()
	if err != nil {
		logger.Warning(err)
		return nil
	}
	return message
}
