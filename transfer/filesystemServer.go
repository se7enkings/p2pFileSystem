package transfer

import (
	"net"
    "fmt"
    "encoding/binary"
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
		go handleConn(conn)
	}
}

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}

// max size is 2GB
func handleConn(conn net.Conn) {
    buff := make([]byte, 4)
    conn.Read(buff)
    connType := string(buff)
    conn.Read(buff)
    binary.LittleEndian.Uint32(buff)
    connSize := int32(buff)
    switch connType {
        case "ndpp":
        fmt.Println(connType, connSize)
    }
	defer conn.Close()
}
