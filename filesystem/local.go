package filesystem

import (
	"github.com/CRVV/p2pFileSystem/settings"
	"io/ioutil"
	"os"
)

func GetLocalFiles(sharedPath string, currentPath string, outputChan chan LocalFile) error {
	fileInfos, err := ioutil.ReadDir(sharedPath)
	if err != nil {
		return err
	}
	for _, v := range fileInfos {
		if v.IsDir() {
			if settings.GetSettings().IsIgnored(v.Name()) {
				continue
			}
			GetLocalFiles(sharedPath+"/"+v.Name(), currentPath+"/"+v.Name(), outputChan)
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
func RemoveLocalFile(fileHash string) {
	FsMutex.Lock()
	file := FileSystem[fileHash]
	FsMutex.Unlock()
	name := settings.GetSettings().GetSharePath() + file.Path + "/" + file.Name
	os.Remove(name)
}
