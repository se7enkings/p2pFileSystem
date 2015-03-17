package transfer

import (
    "github.com/CRVV/p2pFileSystem/settings"
    "net"
    "encoding/binary"
    "fmt"
    "strings"
)

func NeighborSolicitation() {
	message := []byte("CRVV's p2pFileSystem, written in go. " + settings.GetUserName())
	SendNeighborSolicitationMessage("ndpp", message)
}
func SendNeighborSolicitationMessage(header string, message []byte) {
	conn, err := net.Dial("udp", "255.255.255.255:1540")
	checkError(err)
	defer conn.Close()

	conn.Write([]byte("ndpp"))
	buff := make([]byte, 4)
	binary.LittleEndian.PutUint32(buff, uint32(len(message)))
	conn.Write(buff)
	conn.Write(message)
}
func StartNeighborDiscoveryServer() {
	conn, err := net.ListenUDP("udp", &net.UDPAddr{
		IP:   net.IPv4zero,
		Port: 1540})
	checkError(err)
	defer conn.Close()
	for {
		buff := make([]byte, 4)
		conn.Read(buff)
		connType := string(buff)

		conn.Read(buff)
		connSize := binary.LittleEndian.Uint32(buff)

		fmt.Println(connType)
		buff = make([]byte, connSize)
		conn.Read(buff)
        message := strings.Split(string(buff), " ")
        username := message[len(message)-1]
        onReceiveNeighborSolicitation(username, conn.RemoteAddr())
	}
}
func onReceiveNeighborSolicitation(username string, addr net.Addr) {
    fmt.Println(username, addr)
}