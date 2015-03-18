package transfer

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"github.com/CRVV/p2pFileSystem/filesystem"
	"github.com/CRVV/p2pFileSystem/settings"
	"net"
	"time"
    "os"
)

func SendNeighborSolicitation(targetAddr string) {
	conn, err := net.Dial("udp", targetAddr+settings.NeighborDiscoveryPort)
	if err != nil {
		return
	}
	defer conn.Close()

	conn.Write([]byte(settings.NeighborDiscoveryProtocol))
	buff := make([]byte, settings.MessageHeaderSize)

	addr, _, err := net.SplitHostPort(conn.LocalAddr().String())
	if err != nil {
		return
	}
	message, err := ClientMessage2Json(NDMessage{
		Hello:    settings.HelloMessage,
		Username: settings.GetSettings().GetUsername(),
		ID:       filesystem.ID,
		Group:    settings.GetSettings().GetGroupName(),
		Addr:     addr})
	if err != nil {
		return
	}
	binary.LittleEndian.PutUint32(buff, uint32(len(message)))
	conn.Write(buff)
	conn.Write(message)
}
func StartNeighborDiscoveryServer() {
	fmt.Println("start nd server")
	conn, err := net.ListenUDP("udp", &net.UDPAddr{
		IP:   net.IPv4zero,
		Port: 1540})
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	for {
		buff := make([]byte, settings.NeighborDiscoveryMessageBufferSize)
        _, _, err := conn.ReadFromUDP(buff)
        if err != nil {
            continue
        }
        os.Stdout.Write(buff)
        headerSize := settings.MessageHeaderSize
		messageType := string(buff[0:headerSize])
        fmt.Printf("\n\n%s\n", messageType)
//		size := binary.LittleEndian.Uint32(buff[headerSize:headerSize*2])
		if messageType == settings.NeighborDiscoveryProtocol {
			client, err := Json2ClientMessage(buff[headerSize*2:])
			if err != nil {
				continue
			}
			fmt.Println(client)
			onReceiveNeighborSolicitation(client)
		}
	}
}
func onReceiveNeighborSolicitation(client NDMessage) {
	fmt.Println("receive")
	if client.Group == settings.GetSettings().GetGroupName() && client.ID != filesystem.ID {
		fmt.Println("same group, different id")
		_, ok := filesystem.Clients[client.Username]
		if ok || client.Username == settings.GetSettings().GetUsername() {
			fmt.Println("send invalid username")
			SendMessage(client.Addr, settings.InvalidUsername, []byte(settings.InvalidUsername))
		} else {
			fmt.Println("send ndp to received client")
			SendNeighborSolicitation(client.Addr)
			filesystem.OnDiscoverClient(client.Username, client.Addr)
		}
	}
}
func InitNeighborDiscovery() {
	go StartNeighborDiscoveryServer()
	timer := time.Tick(time.Second * 5)
	for {
		SendNeighborSolicitation(settings.BroadcastAddress)
		<-timer
	}
}

type NDMessage struct {
	Hello    string
	Username string
	ID       string
	Group    string
	Addr     string
}

func ClientMessage2Json(client NDMessage) ([]byte, error) {
	b, err := json.Marshal(client)
	return b, err
}
func Json2ClientMessage(jsonClientMessage []byte) (NDMessage, error) {
	cm := NDMessage{}
	err := json.Unmarshal(jsonClientMessage, &cm)
	return cm, err
}
