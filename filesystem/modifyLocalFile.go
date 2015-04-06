package filesystem

import (
	"github.com/CRVV/p2pFileSystem/settings"
	"os"
)

func RemoveLocalFile(fileHash string) {
	filesystemRemote.RLock()
	file := filesystemRemote.M[fileHash]
	filesystemRemote.RUnlock()
	name := settings.GetSettings().GetSharePath() + file.Path + "/" + file.Name
	os.Remove(name)
	RefreshLocalFile()
}
func RemoveDir(path string) {

}
func MakeDir(path string) {

}
func Rename(path0 string, path1 string) {

}
