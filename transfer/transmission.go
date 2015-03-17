package transfer

import (
	"encoding/binary"
	"net"
)

func SendMessage(header string, message []byte) {
	conn, err := net.Dial("tcp", "192.168.10.138:1539")
	checkError(err)

	defer conn.Close()
	conn.Write([]byte("ndpp"))
	buff := make([]byte, 4)
	binary.LittleEndian.PutUint32(buff, uint32(len(message)))
	conn.Write(buff)
	conn.Write(message)
}
