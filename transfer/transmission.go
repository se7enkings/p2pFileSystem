package transfer

import (
	"encoding/binary"
	"github.com/CRVV/p2pFileSystem/filesystem"
	"github.com/CRVV/p2pFileSystem/logger"
	"github.com/CRVV/p2pFileSystem/settings"
	"net"
)

var messagePipe chan Message = make(chan Message, 4)

func sendMessage(addr string, messageType string, message []byte) {
	conn, err := net.Dial("tcp", addr+settings.CommunicationPort)
	if err != nil {
		return err
	}
	defer conn.Close()
	mType := []byte(messageType)
	if len(mType) != settings.MessageHeaderSize {
		logger.Warning("invalid message type")
		return
	}
	conn.Write(mType)

	buff := make([]byte, settings.MessageHeaderSize)
	size := uint32(len(message))
	if size > settings.MaxMessageSize {
		logger.Warning("the message is too long")
		return
	}
	binary.LittleEndian.PutUint32(buff, size)
	conn.Write(buff)
	conn.Write(message)
	logger.Info("sent a " + messageType + " message to " + addr)
}

func FindMessageAndSend() {
	for {
		select {
		case messageFromFS := <-filesystem.MessagePipe:
			sendMessage(clients[messageFromFS.DestinationUsername].Addr, messageFromFS.Type, messageFromFS.Load)
		case message := <-messagePipe:
			sendMessage(message.Destination, message.Type, message.Load)
		}
	}
}

type Message struct {
	Type        string
	Destination string
	Load        []byte
}
