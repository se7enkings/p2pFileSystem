package main

import (
	"github.com/CRVV/p2pFileSystem/filesystem"
	"github.com/CRVV/p2pFileSystem/ndp"
	"github.com/CRVV/p2pFileSystem/remote"
	"github.com/CRVV/p2pFileSystem/ui"
	"runtime"
)

func main() {
	runtime.GOMAXPROCS(8)
	defer ndp.OnExit()

	filesystem.RefreshLocalFile()
	ndp.Init()

	go remote.MaintainClientList()
	go ndp.StartNeighborDiscoveryServer()

	go remote.StartFilesystemServer()

	go ui.StartHttpServer()
	ui.StartCLI()
}
