package filesystem

type File struct {
    Name string // filename, exclude path
    Path string // path in P2P filesystem, "/movie"
//    LocalPath string // path in OS filesystem, "C:/User" or "/home"
    Size int64 // bytes
    FileHash [32]byte // SHA-256 Hash
//    BlockHash [][]byte // SHA-256 Hash
//    Owner User
//    Permission byte
}

//type Directory struct {
//    Name string // directory name, exclude path
//    Files []*File // files contained in this directory
//    SubDirectories []*Directory // directories contained in this directory
//    Owner User
//    Permission byte
//}

//type FileSystemNode interface {
//
//}

//type User struct {
//    Name string
//}
//
//type Group struct {
//    Name string
//    Member []User
//}
