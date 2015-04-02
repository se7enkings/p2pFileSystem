package transfer

import (
	"encoding/binary"
	"fmt"
	"github.com/CRVV/p2pFileSystem/logger"
	"github.com/CRVV/p2pFileSystem/settings"
	"io"
	"net"
)

func SendTcpMessage(message Message) error {
	_, err := sendMessage(message)
	return err
}
func TcpConnectionForReceiveFile(message Message) ([]byte, error) {
	return sendMessage(message)
}
func sendMessage(message Message) ([]byte, error) {
	addr := message.Destination()
	messageType := message.Type()
	payload := message.Payload()

	conn, err := net.Dial("tcp", addr+settings.CommunicationPort)
	if err != nil {
		logger.Warning(err)
		return nil, err
	}
	defer conn.Close()
	mType := []byte(messageType)
	if len(mType) != settings.MessageHeaderSize {
		logger.Warning("invalid message type")
		return nil, err
	}
	conn.Write(mType)

	buff := make([]byte, settings.MessageHeaderSize)
	size := uint32(len(payload))
	if size > settings.MaxMessageSize {
		logger.Warning("the message is too long")
		return nil, err
	}
	binary.LittleEndian.PutUint32(buff, size)
	conn.Write(buff)
	conn.Write(payload)
	logger.Info("sent a " + messageType + " message to " + addr)

	if message.Type() == settings.FileBlockRequestProtocol {
		buffSize := settings.FileBlockSize
		buff = make([]byte, buffSize)
		size, err := io.ReadFull(conn, buff)
		if err != nil && err != io.ErrUnexpectedEOF {
			logger.Warning(err)
			return nil, err
		}
		logger.Info(fmt.Sprintf("receive %d bytes for this block", size))
		return buff[:size], nil
	}
	return nil, nil
}
