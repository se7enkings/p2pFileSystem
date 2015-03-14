package local

import (
	"github.com/CRVV/p2pFileSystem/settings"
	"io/ioutil"
)

func ReadFiles(sharedPath string, currentPath string, outputChan chan LocalFile) error {
	fileInfos, err := ioutil.ReadDir(sharedPath)
	if err != nil {
		return err
	}
	for _, v := range fileInfos {
		if v.IsDir() {
			if settings.IsIgnoredDir(v.Name()) {
				continue
			}
//			outputChan <- LocalFile{currentPath, v}
			ReadFiles(sharedPath+"/"+v.Name(), currentPath+"/"+v.Name(), outputChan)
		} else {
			if !settings.IsIgnoredFile(v.Name()) {
				outputChan <- LocalFile{currentPath, v}
			}
		}
	}
	if currentPath == "" {
		close(outputChan)
	}
	return nil
}
