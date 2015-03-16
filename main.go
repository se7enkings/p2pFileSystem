package main

import (
	"fmt"
	"github.com/CRVV/p2pFileSystem/filesystem"
	"os"
	"runtime"
    "github.com/CRVV/p2pFileSystem/transfer"
    "github.com/CRVV/p2pFileSystem/settings"
)

func main() {
	runtime.GOMAXPROCS(4)

	fileSystem, err := filesystem.ReadLocalFile(settings.GetSharePath())
	checkError(err)

    _, err = filesystem.GetFileList(fileSystem)
	checkError(err)

    remoteFileSystem, err := filesystem.ReadLocalFile("test/testRemoteFolder")
    checkError(err)
    jsonFileSystemMessage, err := transfer.FileSystem2Json(remoteFileSystem)
    checkError(err)

    remoteFileSystem, err = transfer.Json2FileSystem(jsonFileSystemMessage)
    checkError(err)
    _, err = filesystem.GetFileList(remoteFileSystem)
    checkError(err)
}

func checkError(err error) {
	if err != nil {
		fmt.Println("error: ", err)
		os.Exit(1)
	}
}
