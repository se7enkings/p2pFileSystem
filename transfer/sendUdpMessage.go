package transfer

import (
	"encoding/binary"
	"github.com/CRVV/p2pFileSystem/logger"
	"github.com/CRVV/p2pFileSystem/settings"
	"net"
)

func SendUdpPackage(message Message) error {
	addr := message.Destination()
	messageType := message.Type()
	payload := message.Payload()

	udpAddr, err := net.ResolveUDPAddr("udp", addr+settings.NeighborDiscoveryPort)
	if err != nil {
		logger.Warning(err)
		return err
	}
	conn, err := net.DialUDP("udp", nil, udpAddr)
	if err != nil {
		logger.Warning(err)
		return err
	}
	defer conn.Close()
	var buff []byte
	buff = append(buff, []byte(messageType)...)

	size := make([]byte, settings.MessageHeaderSize)
	binary.LittleEndian.PutUint32(size, uint32(len(payload)))
	buff = append(buff, size...)
	buff = append(buff, payload...)
	_, err = conn.Write(buff)
	if err != nil {
		logger.Warning(err)
		return err
	}
	logger.Info("a udp package has been sent to " + addr)
	return nil
}
