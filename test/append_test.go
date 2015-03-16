package test

import (
	"github.com/CRVV/p2pFileSystem/filesystem"
	"github.com/CRVV/p2pFileSystem/transfer"
	"reflect"
	"testing"
)

// the test function is wrong
func TestAppendFilesystem(t *testing.T) {
	fileSystemLocal, _ := filesystem.ReadLocalFile(LocalFolder)
	fileSystemRemote, _ := filesystem.ReadLocalFile(RemoteFolder)
	fileSystemResult, _ := filesystem.ReadLocalFile(ResultFolder)

	jsonMessageRemote, _ := transfer.FileSystem2Json(fileSystemRemote)
	receivedFileSystem, _ := transfer.Json2FileSystem(jsonMessageRemote)

	fileSystemAppended := filesystem.AppendFilesystem(fileSystemLocal, receivedFileSystem)
	if !reflect.DeepEqual(fileSystemResult, fileSystemAppended) {
		t.Errorf("error: \n %v \n %v", fileSystemAppended, fileSystemResult)
	}
}

func AppendFilesystem(originFileSystem filesystem.Filesystem, receivedFileSystem filesystem.Filesystem) filesystem.Filesystem {
	return filesystem.AppendFilesystem(originFileSystem, receivedFileSystem)
}
