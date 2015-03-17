package transfer
import (
    "net"
    "encoding/binary"
)

func SendMessage(header string, message []byte) {
    conn, err := net.Dial("tcp", "localhost:1539")
    defer conn.Close()
    checkError(err)
    conn.Write([]byte("ndpp"))
    buff := make([]byte, 4)
    binary.LittleEndian.PutUint32(buff, uint32(len(message)))
    conn.Write(buff)
    conn.Write(message)
}
