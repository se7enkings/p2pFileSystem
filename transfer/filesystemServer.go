package transfer

import (
	"net"
    "fmt"
    "encoding/binary"
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
		go handleConn(conn)
	}
}


// max size is 4GB
func handleConn(conn net.Conn) {
    buff := make([]byte, 4)
    conn.Read(buff)
    connType := string(buff)

    conn.Read(buff)
    connSize := binary.LittleEndian.Uint32(buff)

    switch connType {
        case "ndpp":
        fmt.Println(connType)
        message := make([]byte, connSize)
        conn.Read(message)
        messageStr := strings.Split(string(message), " ")
        fmt.Println(messageStr[len(messageStr)-1])
    }
	defer conn.Close()
}
func checkError(err error) {
    if err != nil {
        panic(err)
    }
}
