package gui

import (
    "github.com/CRVV/p2pFileSystem/filesystem"
    "strings"
    "fmt"
)

var fileList Node = Node{"root", true, 0, *new([32]byte), make(map[string]Node)}

func GetFileList(fileSystem map[[32]byte]filesystem.File) (Node, error) {
    for fileHash, file := range fileSystem{
        folder := createFolder(fileList, file.Path)
        folder.Children[file.Name] = Node{file.Name, false, file.Size, fileHash, nil}
    }
    fmt.Println(fileList)
    return fileList, nil
}

func StartGuiServer(fileList Node) error {
    return nil
}

func createFolder(rootFolder Node, folder string) Node {
    folders := strings.Split(folder, "/")
    return doCreateFolder(rootFolder, folders)
}

func doCreateFolder(rootFolder Node, folders []string) Node {
    _, ok := rootFolder.Children[folders[0]]
    if !ok && folders[0] != ""{
        rootFolder.Children[folders[0]] = Node{folders[0], true, 0, *new([32]byte), make(map[string]Node)}
    }
    if len(folders) > 1 {
        return doCreateFolder(rootFolder.Children[folders[0]], folders[1:])
    }
    return rootFolder.Children[folders[0]]
}
