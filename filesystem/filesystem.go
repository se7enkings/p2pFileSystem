package filesystem

import (
	"crypto/sha256"
	"encoding/base64"
	"github.com/CRVV/p2pFileSystem/logger"
	"github.com/CRVV/p2pFileSystem/settings"
	"io"
	"os"
	"strings"
	"sync"
)

var FileSystem Filesystem
var FsMutex sync.Mutex = sync.Mutex{}

var FileList Node
var FlMutex sync.Mutex = sync.Mutex{}

func readLocalFile(folder string) error {
	fileSystem := make(Filesystem)
	filesChan := make(chan LocalFile)

	go getLocalFiles(folder, "", filesChan)

	for f := range filesChan {
		sha256Sum, err := getFileHash(folder + "/" + f.Path + "/" + f.FileInfo.Name())
		if err != nil {
			return err
		}
		hash := base64.URLEncoding.EncodeToString(sha256Sum[:])
		fileSystem[hash] = &File{
			Name:    f.FileInfo.Name(),
			Path:    f.Path,
			Size:    f.FileInfo.Size(),
			AtLocal: true,
			Owner:   []string{settings.GetSettings().GetUsername()}}
	}
	FsMutex.Lock()
	FileSystem = fileSystem
	FsMutex.Unlock()
	FslMutex.Lock()
	FileSystemLocal = fileSystem
	FslMutex.Unlock()
	return nil
}
func getFileHash(name string) ([]byte, error) {
	file, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	hasher := sha256.New()
	_, err = io.Copy(hasher, file)
	if err != nil {
		return nil, err
	}
	return hasher.Sum(nil), nil
}

func generateFileList() error {
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
	err := readLocalFile(settings.GetSettings().GetSharePath())
	logger.Error(err)
	CMutex.Lock()
	if Clients == nil {
		Clients = make(map[string]Client)
	}
	for username, c := range Clients {
		FsMutex.Lock()
		appendFilesystem(FileSystem, c.FileSystem, username)
		FsMutex.Unlock()
	}
	CMutex.Unlock()
	err = generateFileList()
	logger.Error(err)
}
func appendFilesystem(originFileSystem Filesystem, receivedFileSystem Filesystem, username string) {
	// because the File.IsLocal is ignored by json, IsLocal in received Filesystem is always default bool value(false).
	// It is possible that duplicate filename exists in the returned Filesystem(the content is different).
	logger.Info("appedn filesystem begin")
	for hash, file := range receivedFileSystem {
		_, ok := originFileSystem[hash]
		if ok {
			originFileSystem[hash].Owner = append(originFileSystem[hash].Owner, username)
		} else {
			originFileSystem[hash] = file
		}
	}
	logger.Info("appedn filesystem end")
}
