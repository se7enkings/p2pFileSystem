package filesystem

import (
	"github.com/CRVV/p2pFileSystem/logger"
	"github.com/CRVV/p2pFileSystem/settings"
	"os"
)

func RemoveLocalFile(fileHash string) {
	filesystemLocal.RLock()
	file, _ := filesystemLocal.M[fileHash]
	filesystemLocal.RUnlock()
	name := settings.GetSettings().GetSharePath() + file.Path + "/" + file.Name
	err := os.Remove(name)
	logger.Warning(err)
	RefreshLocalFile()
}
func RemoveDir(path string) {

}
func MakeDir(path string) {

}
func Rename(path0 string, path1 string) {

}
