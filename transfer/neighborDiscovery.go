package transfer

import "github.com/CRVV/p2pFileSystem/settings"

func NeighborSolicitation() {
    message := []byte("CRVV's p2pFileSystem, written in go. "+settings.GetUserName())
    SendMessage("ndpp", message)
}
