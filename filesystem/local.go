package filesystem

import (
	"crypto/sha256"
	"encoding/base64"
	"github.com/CRVV/p2pFileSystem/settings"
	"io"
	"io/ioutil"
	"os"
)

func readLocalFile() (map[string]*File, error) {
	folder := settings.GetSettings().GetSharePath()
	fileSystemLocalTemp := make(map[string]*File)
	filesChan := make(chan LocalFile, 8)

	go readFolder(folder, "", filesChan)

	for f := range filesChan {
		sha256Sum, err := getFileHash(folder + "/" + f.Path + "/" + f.FileInfo.Name())
		if err != nil {
			return nil, err
		}
		hash := base64.URLEncoding.EncodeToString(sha256Sum[:])
		fileSystemLocalTemp[hash] = &File{
			Name:    f.FileInfo.Name(),
			Path:    f.Path,
			Size:    f.FileInfo.Size(),
			AtLocal: true,
		}
	}
	return fileSystemLocalTemp, nil
}
func getFileHash(name string) ([]byte, error) {
	file, err := os.Open(name)
	defer file.Close()
	if err != nil {
		return nil, err
	}
	hasher := sha256.New()
	_, err = io.Copy(hasher, file)
	if err != nil {
		//TODO: should handle this error
		return nil, err
	}
	return hasher.Sum(nil), nil
}

func readFolder(sharedPath string, currentPath string, outputChan chan LocalFile) error {
	fileInfos, err := ioutil.ReadDir(sharedPath)
	if err != nil {
		return err
	}
	for _, v := range fileInfos {
        if v.Mode() == os.ModeSymlink {
            continue
        }
		if v.IsDir() {
			if settings.GetSettings().IsIgnored(v.Name()) {
				continue
			}
			readFolder(sharedPath+"/"+v.Name(), currentPath+"/"+v.Name(), outputChan)
		} else {
			if !settings.GetSettings().IsIgnored(v.Name()) {
				outputChan <- LocalFile{currentPath, v}
			}
		}
	}
	if currentPath == "" {
		close(outputChan)
	}
	return nil
}
