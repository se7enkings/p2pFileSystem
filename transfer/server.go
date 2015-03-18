package transfer

import (
	"encoding/binary"
	"fmt"
	"net"
	"strings"
)

func StartFilesystemServer() {
	// TODO: use const port in settings
	//	addr, err := net.ResolveTCPAddr("tcp", ":1539")
	//	checkError(err)
	listener, err := net.Listen("tcp", ":1539")
	checkError(err)
	defer listener.Close()
	for {
		conn, err := listener.Accept()
		checkError(err)
		go handleTcpConn(conn)
	}
}

// max size is 4GB
func handleTcpConn(conn net.Conn) {
	buff := make([]byte, 4)
	conn.Read(buff)
	connType := string(buff)

	conn.Read(buff)
	connSize := binary.LittleEndian.Uint32(buff)

	switch connType {
	case "ndpp":
		fmt.Println(connType)
		buff = make([]byte, connSize)
		conn.Read(buff)
		message := strings.Split(string(buff), " ")
		fmt.Println(message[len(message)-1])
	}
	defer conn.Close()
}
func checkError(err error) {
	if err != nil {
		panic(err)
	}
}
