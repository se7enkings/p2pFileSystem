package transfer

import (
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/CRVV/p2pFileSystem/filesystem"
	"github.com/CRVV/p2pFileSystem/logger"
	"github.com/CRVV/p2pFileSystem/settings"
	"net"
)

func StartFilesystemServer() {
	logger.Info("start filesystem server")
	listener, err := net.Listen("tcp", settings.CommunicationPort)
	if err != nil {
		logger.Error(err)
	}
	defer listener.Close()
	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		go handleTcpConn(conn)
	}
}

// max size is 4GB
func handleTcpConn(conn net.Conn) {
	logger.Info("handle a tcp connection from " + conn.RemoteAddr().String())
	defer conn.Close()
	buff := make([]byte, settings.MessageHeaderSize)
	conn.Read(buff)
	messageType := string(buff)

	conn.Read(buff)
	size := binary.LittleEndian.Uint32(buff)

	buff = make([]byte, size)
	conn.Read(buff)

	switch messageType {
	case settings.FileSystemListProtocol:
		logger.Info("receive a fileSystemList message from " + conn.RemoteAddr().String())
		filesystem.OnReceiveFilesystem(buff)
	case settings.InvalidUsername:
		logger.Info("receive a invalidUsername message from " + conn.RemoteAddr().String())
		onReceiveInvalidUsername()
	}
}

func onReceiveInvalidUsername() {
	logger.Error(errors.New("duplicate username"))
}
