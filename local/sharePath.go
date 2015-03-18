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
			if settings.GetSettings().IsIgnored(v.Name()) {
				continue
			}
			ReadFiles(sharedPath+"/"+v.Name(), currentPath+"/"+v.Name(), outputChan)
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
