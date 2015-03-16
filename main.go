package main

import (
	"fmt"
	"github.com/CRVV/p2pFileSystem/filesystem"
	"github.com/CRVV/p2pFileSystem/gui"
	"os"
)

func main() {
	fileSystem, err := filesystem.ReadLocalFile()
	errorChecker(err)

	guiFileList, err := gui.GetFileList(fileSystem)
	errorChecker(err)

    fmt.Println(guiFileList)

	go gui.StartGuiServer(guiFileList)
}

func errorChecker(err error) {
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}
