package main

import (
	"github.com/CRVV/p2pFileSystem/filesystem"
	"github.com/CRVV/p2pFileSystem/transfer"
	"github.com/CRVV/p2pFileSystem/ui"
	"runtime"
)

func main() {
	runtime.GOMAXPROCS(4)
	filesystem.Init()
	go transfer.InitNeighborDiscovery()
	go transfer.StartFilesystemServer()
	//	transfer.SendNeighborSolicitation("192.168.10.115")
	ui.StartCLI()
}
