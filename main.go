package main

import (
	"fmt"
	"github.com/CRVV/p2pFileSystem/filesystem"
	"os"
	"runtime"
    "github.com/CRVV/p2pFileSystem/gui"
    "github.com/CRVV/p2pFileSystem/transfer"
)

func main() {
	runtime.GOMAXPROCS(4)

	fileSystem, err := filesystem.ReadLocalFile()
	checkError(err)

	fileList, err := filesystem.GetFileList(fileSystem)
	checkError(err)

	go gui.StartGuiServer(fileList)

    jsonFileListMessage, err := transfer.FileSystem2Json(fileSystem)
    checkError(err)

    fileSystem2, err := transfer.Json2FileSystem(jsonFileListMessage)
    checkError(err)
    os.Stdout.Write(jsonFileListMessage)
    fmt.Println()
    fmt.Println(fileSystem2)
}

func checkError(err error) {
	if err != nil {
		fmt.Println("error: ", err)
		os.Exit(1)
	}
}
