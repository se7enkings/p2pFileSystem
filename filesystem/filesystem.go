package filesystem

import (
	"crypto/sha256"
	"encoding/base64"
	"github.com/CRVV/p2pFileSystem/settings"
	"io/ioutil"
	"strings"
)

var FileSystem Filesystem
var FileSystemLocal Filesystem
var FileList Node

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
	FileSystem = fileSystem
	FileSystemLocal = fileSystem
	return nil
}

func GetFileList() error {
	fileList := Node{"root", true, true, 0, "", make(map[string]*Node)}
	fileList.Children[".."] = &fileList
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
	FileList = fileList
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
	checkError(err)
	for _, c := range Clients {
		FileSystem = AppendFilesystem(FileSystem, c.FileSystem)
	}
	err = GetFileList()
	checkError(err)
}
func checkError(err error) {
	if err != nil {
		panic(err)
	}
}
