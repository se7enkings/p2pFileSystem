package local
import (
    "os"
    "fmt"
)

type LocalFile struct {
    Path string
    FileInfo os.FileInfo
}

func (localFile LocalFile) String() string{
    return fmt.Sprintf("%v, %v", localFile.Path, localFile.FileInfo.Name())
}