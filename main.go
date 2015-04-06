package main

import (
	"github.com/CRVV/p2pFileSystem/filesystem"
	"github.com/CRVV/p2pFileSystem/ndp"
	"github.com/CRVV/p2pFileSystem/ui"
	"runtime"
)

func main() {
	runtime.GOMAXPROCS(8)
	defer ndp.OnExit()

	filesystem.Init()
	ndp.Init()

	go filesystem.MaintainClientList()
	go ndp.StartNeighborDiscoveryServer()

	go filesystem.StartFilesystemServer()

	ui.StartCLI()
}
