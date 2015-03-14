package filesystem
import (
    "github.com/CRVV/p2pFileSystem/local"
    "github.com/CRVV/p2pFileSystem/settings"
    "fmt"
    "crypto/sha256"
    "io/ioutil"
)

//var RootDirectory Directory

func ReadLocalFile() ([]File, error) {
    var files []File
    filesChan := make(chan local.LocalFile)
    go local.ReadFiles(settings.GetSharePath(), "", filesChan)

    for f := range filesChan {
//        fmt.Println(f)
        fileData, err := ioutil.ReadFile(settings.GetSharePath()+"/"+f.Path+"/"+f.FileInfo.Name())
        if err != nil{
            return nil, err
        }
        sha256Sum := sha256.Sum256(fileData)
        files = append(files, File{f.FileInfo.Name(), f.Path, f.FileInfo.Size(), sha256Sum})
    }
    fmt.Println(files)
    return files, nil
}
