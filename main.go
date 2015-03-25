package main

import (
	"github.com/CRVV/p2pFileSystem/filesystem"
	"github.com/CRVV/p2pFileSystem/transfer"
	"github.com/CRVV/p2pFileSystem/ui"
	"runtime"
)

func main() {
	runtime.GOMAXPROCS(4)
	defer transfer.OnExit()

	filesystem.Init()
	go transfer.InitNeighborDiscovery()
	go transfer.StartFilesystemServer()
	go transfer.FindMessageAndSend()
	go transfer.NeighborDiscovery()
	ui.StartCLI()
}
