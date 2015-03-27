package main

import (
	"github.com/CRVV/p2pFileSystem/filesystem"
	"github.com/CRVV/p2pFileSystem/ndp"
	"github.com/CRVV/p2pFileSystem/ui"
	"runtime"
)

func main() {
	runtime.GOMAXPROCS(4)
	defer ndp.OnExit()

	filesystem.Init()

	go ndp.StartNeighborDiscoveryServer()
	go filesystem.MaintainClientList()

	go filesystem.StartFilesystemServer()

	ui.StartCLI()
}
