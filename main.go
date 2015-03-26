package main

import (
	"github.com/CRVV/p2pFileSystem/filesystem"
	"github.com/CRVV/p2pFileSystem/ndp"
	"github.com/CRVV/p2pFileSystem/transfer"
	"github.com/CRVV/p2pFileSystem/ui"
	"runtime"
)

func main() {
	runtime.GOMAXPROCS(4)
	defer ndp.OnExit()

	filesystem.Init()

	go ndp.StartNeighborDiscoveryServer()
	go ndp.NeighborDiscovery()

	go transfer.StartFilesystemServer()

	ui.StartCLI()
}
