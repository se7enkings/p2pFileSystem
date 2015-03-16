package filesystem

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"github.com/CRVV/p2pFileSystem/local"
	"io/ioutil"
	"strings"
)

//var RootDirectory Directory

func ReadLocalFile(folder string) (Filesystem, error) {
	fileSystem := make(Filesystem)
	filesChan := make(chan local.LocalFile)

	//    go local.ReadFiles(settings.GetSharePath(), "", filesChan)
	go local.ReadFiles(folder, "", filesChan)

	for f := range filesChan {
		// TODO: ioutil.ReadFile cannot handle big file(use too much memory), fix it.
		fileData, err := ioutil.ReadFile(folder + "/" + f.Path + "/" + f.FileInfo.Name())
		if err != nil {
			return nil, err
		}
		sha256Sum := sha256.Sum256(fileData)
		hash := base64.StdEncoding.EncodeToString(sha256Sum[:])
		fileSystem[hash] = File{f.FileInfo.Name(), f.Path, f.FileInfo.Size(), true}
	}
	return fileSystem, nil
}

var fileList Node = Node{"root", true, true, 0, "", make(map[string]Node)}

func GetFileList(fileSystem Filesystem) (Node, error) {
	for fileHash, file := range fileSystem {
		folder := createFolder(fileList, file.Path)
		folder.Children[file.Name] = Node{file.Name, false, file.AtLocal, file.Size, fileHash, nil}
	}
	fmt.Println(fileList)
	return fileList, nil
}
func createFolder(rootFolder Node, folder string) Node {
	folders := strings.Split(folder, "/")
	return doCreateFolder(rootFolder, folders[1:])
}
func doCreateFolder(rootFolder Node, folders []string) Node {
	if len(folders) == 0 {
		return rootFolder
	}
	_, ok := rootFolder.Children[folders[0]]
	if !ok && folders[0] != "" {
		rootFolder.Children[folders[0]] = Node{folders[0], true, true, 0, "", make(map[string]Node)}
	}
	if len(folders) > 1 {
		return doCreateFolder(rootFolder.Children[folders[0]], folders[1:])
	}
	return rootFolder.Children[folders[0]]
}
