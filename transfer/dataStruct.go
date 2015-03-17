package transfer

import (
    "github.com/CRVV/p2pFileSystem/filesystem"
    "encoding/json"
)

func FileSystem2Json(fileSystem filesystem.Filesystem) ([]byte, error){
    b, err := json.Marshal(fileSystem)
    return b, err
}

func Json2FileSystem(jsonFileListMessage []byte) (filesystem.Filesystem, error){
    fs := make(filesystem.Filesystem)
    err := json.Unmarshal(jsonFileListMessage, &fs)
    return fs, err
}
