package filesystem

import (
	"crypto/sha256"
	"encoding/base64"
	"github.com/CRVV/p2pFileSystem/logger"
	"github.com/CRVV/p2pFileSystem/settings"
	"io/ioutil"
	"strings"
	"sync"
)

var FileSystem Filesystem
var FsMutex sync.Mutex = sync.Mutex{}

var FileList Node
var FlMutex sync.Mutex = sync.Mutex{}

func ReadLocalFile(folder string) error {
	fileSystem := make(Filesystem)
	filesChan := make(chan LocalFile)

	go GetLocalFiles(folder, "", filesChan)

	for f := range filesChan {
		// TODO: ioutil.ReadFile cannot handle big file(use too much memory), fix it.
		fileData, err := ioutil.ReadFile(folder + "/" + f.Path + "/" + f.FileInfo.Name())
		if err != nil {
			return err
		}
		sha256Sum := sha256.Sum256(fileData)
		hash := base64.StdEncoding.EncodeToString(sha256Sum[:])
		fileSystem[hash] = File{f.FileInfo.Name(), f.Path, f.FileInfo.Size(), true}
	}
	FsMutex.Lock()
	FileSystem = fileSystem
	FsMutex.Unlock()
	FslMutex.Lock()
	FileSystemLocal = fileSystem
	FslMutex.Unlock()
	return nil
}

func GetFileList() error {
	fileList := Node{"root", true, true, 0, "", make(map[string]*Node)}
	fileList.Children[".."] = &fileList
	FsMutex.Lock()
	for fileHash, file := range FileSystem {
		folder := createFolder(&fileList, file.Path)
		_, ok := folder.Children[file.Name]
		name := file.Name
		if ok {
			name += "-1"
			//TODO: do better on duplicate filename. This will produce filename "xxx.txt-1"
		}
		folder.Children[name] = &Node{name, false, file.AtLocal, file.Size, fileHash, nil}
	}
	FsMutex.Unlock()
	FlMutex.Lock()
	FileList = fileList
	FlMutex.Unlock()
	return nil
}
func createFolder(rootFolder *Node, folder string) *Node {
	folders := strings.Split(folder, "/")
	return doCreateFolder(rootFolder, folders[1:])
}
func doCreateFolder(rootFolder *Node, folders []string) *Node {
	if len(folders) == 0 {
		return rootFolder
	}
	_, ok := rootFolder.Children[folders[0]]
	if !ok && folders[0] != "" {
		rootFolder.Children[folders[0]] = &Node{folders[0], true, true, 0, "", make(map[string]*Node)}
		rootFolder.Children[folders[0]].Children[".."] = rootFolder
	}
	if len(folders) > 1 {
		return doCreateFolder(rootFolder.Children[folders[0]], folders[1:])
	}
	return rootFolder.Children[folders[0]]
}
func Init() {
	err := ReadLocalFile(settings.GetSettings().GetSharePath())
	logger.Error(err)
	CMutex.Lock()
	if Clients == nil {
		Clients = make(map[string]Client)
	}
	for _, c := range Clients {
		FsMutex.Lock()
		FileSystem = AppendFilesystem(FileSystem, c.FileSystem)
		FsMutex.Unlock()
	}
	CMutex.Unlock()
	err = GetFileList()
	logger.Error(err)
}
