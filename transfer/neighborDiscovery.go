package transfer

import (
    "github.com/CRVV/p2pFileSystem/settings"
    "net"
    "encoding/binary"
    "fmt"
)

func NeighborSolicitation() {
    hello := "CRVV's p2pFileSystem, written in go."

    // TODO: move this to settings
	conn, err := net.Dial("udp", "255.255.255.255:1540")
	checkError(err)
	defer conn.Close()

	conn.Write([]byte("ndpp"))
	buff := make([]byte, 4)

    addr, _, err := net.SplitHostPort(conn.LocalAddr().String())
    checkError(err)
    message, err := ClientMessage2Json(Client{hello, settings.GetUserName(), addr})
    checkError(err)

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

		buff = make([]byte, connSize)
		conn.Read(buff)
        if connType == "ndpp" {
            client, err := Json2ClientMessage(buff)
            checkError(err)
            onReceiveNeighborSolicitation(client)
        }
	}
}
func onReceiveNeighborSolicitation(client Client) {
    fmt.Println(client)
}