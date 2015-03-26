package transfer

import (
	"encoding/binary"
	"github.com/CRVV/p2pFileSystem/logger"
	"github.com/CRVV/p2pFileSystem/settings"
	"net"
)

func SendTcpMessage(message Message) error {
	addr := message.Destination()
	messageType := message.Type()
	payload := message.Payload()

	conn, err := net.Dial("tcp", addr+settings.CommunicationPort)
	if err != nil {
		logger.Warning(err)
		return err
	}
	defer conn.Close()
	mType := []byte(messageType)
	if len(mType) != settings.MessageHeaderSize {
		logger.Warning("invalid message type")
		return err
	}
	conn.Write(mType)

	buff := make([]byte, settings.MessageHeaderSize)
	size := uint32(len(payload))
	if size > settings.MaxMessageSize {
		logger.Warning("the message is too long")
		return err
	}
	binary.LittleEndian.PutUint32(buff, size)
	conn.Write(buff)
	conn.Write(payload)
	logger.Info("sent a " + messageType + " message to " + addr)
	return nil
}
