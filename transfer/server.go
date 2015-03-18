package transfer

import (
	"encoding/binary"
	"github.com/CRVV/p2pFileSystem/filesystem"
	"github.com/CRVV/p2pFileSystem/settings"
	"net"
    "errors"
)

func StartFilesystemServer() {
	listener, err := net.Listen("tcp", settings.CommunicationPort)
	if err != nil {
		panic(err)
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
		filesystem.OnReceiveFilesystem(buff)
    case settings.InvalidUsername:
        onReceiveInvalidUsername()
	}
}

func onReceiveInvalidUsername() {
    panic(errors.New("duplicate username"))
}
