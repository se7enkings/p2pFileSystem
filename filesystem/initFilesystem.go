package filesystem
import (
    "github.com/CRVV/p2pFileSystem/local"
    "github.com/CRVV/p2pFileSystem/settings"
//    "fmt"
    "crypto/sha256"
    "io/ioutil"
)

//var RootDirectory Directory

func ReadLocalFile() (map[[32]byte]File, error) {
    fileSystem := make(map[[32]byte]File)
    filesChan := make(chan local.LocalFile)

    go local.ReadFiles(settings.GetSharePath(), "", filesChan)

    for f := range filesChan {
//        fmt.Println(f)
        // TODO: ioutil.ReadFile cannot handle big file(use too much memory), fix it.
        fileData, err := ioutil.ReadFile(settings.GetSharePath()+"/"+f.Path+"/"+f.FileInfo.Name())
        if err != nil{
            return nil, err
        }
        sha256Sum := sha256.Sum256(fileData)
        fileSystem[sha256Sum] = File{f.FileInfo.Name(), f.Path, f.FileInfo.Size()}
    }
//    fmt.Println(fileSystem)

    return fileSystem, nil
}
