package filesystem

import (
	"encoding/json"
	"github.com/CRVV/p2pFileSystem/logger"
	"strings"
)

var noticeUI chan int

func GetFileList() *FileList {
	return &fileList
}
func GetFileListJson(notice chan int) []byte {
	noticeUI = notice
	fileList.RLock()
	b, err := json.Marshal(fileList)
	fileList.RUnlock()
	logger.Warning(err)
	return b
}
func generateFileList() *Node {
	fileListTemp := Node{"root", true, true, 0, "", make(map[string]*Node)}
	//TODO: delete this line for json
	//	fileListTemp.Children[".."] = &fileListTemp
	filesystemLocal.RLock()
	for fileHash, file := range filesystemLocal.M {
		addFileToList(&fileListTemp, fileHash, file)
	}
	filesystemRemote.RLock()
	for fileHash, file := range filesystemRemote.M {
		fileAtLocal, ok := filesystemLocal.M[fileHash]
		if ok && fileAtLocal.Path == file.Path {
			continue
		}
		addFileToList(&fileListTemp, fileHash, file)
	}
	filesystemLocal.RUnlock()
	filesystemRemote.RUnlock()
	if noticeUI != nil {
		noticeUI <- 1
	}
	return &fileListTemp
}
func addFileToList(rootFolder *Node, fileHash string, file *File) {
	folder := createFolder(rootFolder, file.Path)
	name := file.Name
	_, ok := folder.Children[name]
	if ok {
		//TODO: do better on duplicate filename. This will produce filename "xxx.txt-1"
		name += "-1"
	}
	folder.Children[name] = &Node{
		Name:     name,
		IsDir:    false,
		AtLocal:  file.AtLocal,
		Size:     file.Size,
		FileHash: fileHash,
	}
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
		//TODO: delete this for json
		//		rootFolder.Children[folders[0]].Children[".."] = rootFolder
	}
	if len(folders) > 1 {
		return doCreateFolder(rootFolder.Children[folders[0]], folders[1:])
	}
	return rootFolder.Children[folders[0]]
}
