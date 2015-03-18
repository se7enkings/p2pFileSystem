package main

import (
	"github.com/CRVV/p2pFileSystem/filesystem"
	"github.com/CRVV/p2pFileSystem/settings"
	"github.com/CRVV/p2pFileSystem/ui"
	"runtime"
)

func main() {
	runtime.GOMAXPROCS(4)

	err := filesystem.ReadLocalFile(settings.GetSettings().GetSharePath())
	checkError(err)

	err = filesystem.GetFileList()
	checkError(err)

	ui.StartCLI()

	//	remoteFileSystem, err := filesystem.ReadLocalFile("test/testRemoteFolder")
	//	checkError(err)
	//	jsonFileSystemMessage, err := transfer.FileSystem2Json(remoteFileSystem)
	//	checkError(err)
	//
	//	remoteFileSystem, err = transfer.Json2FileSystem(jsonFileSystemMessage)
	//	checkError(err)
	//	_, err = filesystem.GetFileList(remoteFileSystem)
	//	checkError(err)

	//    go transfer.StartFilesystemServer()
	//    go transfer.StartNeighborDiscoveryServer()
	//    transfer.NeighborSolicitation()
	//    c := make(chan int)
	//    <- c

	//    ifaces, err := net.Interfaces()
	//    checkError(err)
	//    for _, iface := range ifaces {
	//        fmt.Println(iface.Addrs())
	//    }
}

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}
