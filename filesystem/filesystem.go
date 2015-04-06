package filesystem

import (
	"errors"
	"github.com/CRVV/p2pFileSystem/logger"
)

var filesystemRemote Filesystem
var filesystemLocal Filesystem
var fileList FileList

func RefreshLocalFile() {
	localFiles, err := readLocalFile()
	logger.Error(err)

	filesystemLocal.Lock()
	filesystemLocal.M = localFiles
	filesystemLocal.Unlock()

	fileList.Lock()
	fileList.N = generateFileList()
	fileList.Unlock()
}
func GetRemoteFile(hash string) (name string, path string, size int64, owners []string, err error) {
	filesystemRemote.RLock()
	defer filesystemRemote.RUnlock()
	file, ok := filesystemRemote.M[hash]
	if !ok {
		err = errors.New("no such file")
		return
	}
	name = file.Name
	path = file.Path
	size = file.Size
	for owner, _ := range file.Owners {
		owners = append(owners, owner)
	}
	return
}
func GetLocalFile(hash string) (name string, path string, size int64, owners []string, err error) {
	filesystemLocal.RLock()
	defer filesystemLocal.RUnlock()
	file, ok := filesystemLocal.M[hash]
	if !ok {
		err = errors.New("no such file")
		return
	}
	name = file.Name
	path = file.Path
	size = file.Size
	for owner, _ := range file.Owners {
		owners = append(owners, owner)
	}
	return
}
func RefreshRemoteFile(clients map[string]Filesystem) {
	filesystemRemote.Lock()
	filesystemRemote.M = make(map[string]*File)
	for username, fs := range clients {
		appendFilesystem(filesystemRemote, fs, username)
	}
	filesystemRemote.Unlock()
	fileList.Lock()
	fileList.N = generateFileList()
	fileList.Unlock()
}
func GetLocalFilesystemForSend() *Filesystem {
	return &filesystemLocal
}
func appendFilesystem(originFileSystem Filesystem, receivedFileSystem Filesystem, username string) {
	// because the File.IsLocal is ignored by json, IsLocal in received Filesystem is always default bool value(false).
	// It is possible that duplicate filename exists in the returned Filesystem(the content is different).
	for hash, file := range receivedFileSystem.M {
		_, ok := originFileSystem.M[hash]
		if !ok {
			originFileSystem.M[hash] = file
			originFileSystem.M[hash].Owners = make(map[string]bool)
		}
		// this if is unnecessary
		if !file.AtLocal {
			originFileSystem.M[hash].Owners[username] = true
		}
	}
}
