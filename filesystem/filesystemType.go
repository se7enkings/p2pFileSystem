package filesystem

import (
	"fmt"
	"strings"
)

type Filesystem map[string]File

type File struct {
	Name string // filename, exclude path
	Path string // path in P2P filesystem, "/movie", must have the first "/"
	//    LocalPath string // path in OS filesystem, "C:/User" or "/home"
	Size    int64 // bytes
	AtLocal bool  `json:"-"`
	//    FileHash [32]byte // SHA-256 Hash
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

type Node struct {
	Name     string
	IsDir    bool
	AtLocal  bool
	Size     int64
	FileHash string
	Children map[string]Node
}

func (node Node) String() string {
	return Node2str(&node, 0, true)
}
func Node2str(node *Node, space int, tree bool) string {
	const spaceString string = "    "
	str := ""
	str += strings.Repeat(spaceString, space) + fmt.Sprintf("%s, %s, %s, %d \n", node.Name, node.isDir(), node.atLocal(), node.Size)
	for _, file := range node.Children {
		if file.IsDir && tree {
			str += Node2str(&file, space+1, tree)
		} else {
			str += strings.Repeat(spaceString, space+1) + fmt.Sprintf("%s, %s, %s, %d \n", file.Name, file.isDir(), file.atLocal(), file.Size)
		}
	}
	return str
}
func (node *Node) isDir() string {
	if node.IsDir {
		return "dir"
	} else {
		return "file"
	}
}
func (node *Node) atLocal() string {
	if node.AtLocal {
		return "local"
	} else {
		return "remote"
	}
}
