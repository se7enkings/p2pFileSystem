package main

import (
	"fmt"
	"github.com/CRVV/p2pFileSystem/filesystem"
	"github.com/CRVV/p2pFileSystem/settings"
	"runtime"
)

func main() {
	runtime.GOMAXPROCS(4)

	fileSystem, err := filesystem.ReadLocalFile(settings.GetSettings().GetSharePath())
	checkError(err)

	fileList, err := filesystem.GetFileList(fileSystem)
	checkError(err)

	fmt.Println(fileList)

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
