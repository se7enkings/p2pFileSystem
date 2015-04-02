package filesystem

import (
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/CRVV/p2pFileSystem/logger"
	"github.com/CRVV/p2pFileSystem/ndp"
	"github.com/CRVV/p2pFileSystem/settings"
	"io"
	"net"
)

func StartFilesystemServer() {
	listener, err := net.Listen("tcp", settings.CommunicationPort)
	if err != nil {
		logger.Error(err)
	}
	//	logger.Info("start filesystem server")
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
	messageSize, err := io.ReadFull(conn, buff)
	logger.Warning(err)
	logger.Info(fmt.Sprintf("%d bytes read", messageSize))

	switch messageType {
	case settings.FileSystemListProtocol:
		logger.Info("receive a fileSystemList message from " + conn.RemoteAddr().String())
		OnReceiveFilesystem(buff)
	case settings.FileListRequestProtocol:
		peer, err := ndp.GetPeerFromJson(buff)
		if err != nil {
			logger.Warning(err)
			return
		}
		if peer.Group == settings.GetSettings().GetGroupName() {
			logger.Info("receive a fileSystemRequest message from " + conn.RemoteAddr().String())
			//TODO
			ndp.OnReceiveNeighborSolicitationEcho(peer)
			OnRequestedFilesystem(peer.Username)
		}
	case settings.FileBlockRequestProtocol:
		logger.Info("receive a fileBlockRequest message from " + conn.RemoteAddr().String())
		requestMessage, err := Json2FBRMessage(buff)
		if err != nil {
			logger.Warning(err)
			return
		}
		fileData := onRequestedFileBlock(&requestMessage)
		_, err = conn.Write(fileData)
		logger.Warning(err)
		logger.Info("send file block complete")
	case settings.InvalidUsername:
		logger.Info("receive a invalidUsername message from " + conn.RemoteAddr().String())
		onReceiveInvalidUsername()
	}
}
func onReceiveInvalidUsername() {
	logger.Error(errors.New("duplicate username"))
}
