package filesystem

import (
	"fmt"
	"os"
	"strings"
)

type Filesystem map[string]*File

type File struct {
	Name string // filename, exclude path
	Path string // path in P2P filesystem, "/movie", must have the first "/"
	//    LocalPath string // path in OS filesystem, "C:/User" or "/home"
	Size    int64 // bytes
	AtLocal bool  `json:"-"`
	//    BlockHash [][]byte // SHA-256 Hash
	Owner []string
	//	    Permission byte
}

type Node struct {
	Name     string
	IsDir    bool
	AtLocal  bool
	Size     int64
	FileHash string
	Children map[string]*Node
}

func (node Node) String() string {
	return Node2str(&node, 0, true)
}
func Node2str(node *Node, space int, tree bool) string {
	const spaceString string = "    "
	str := ""
	str += strings.Repeat(spaceString, space) + fmt.Sprintf("%s, %s\n", node.Name, node.isDir())
	for name, file := range node.Children {
		switch {
		case name == "..":
		case tree && file.IsDir:
			str += Node2str(file, space+1, tree)
		case file.IsDir:
			str += strings.Repeat(spaceString, space+1) + fmt.Sprintf("%s, %s\n", file.Name, file.isDir())
		default:
			str += strings.Repeat(spaceString, space+1) + fmt.Sprintf("%s, %s, %s, %d bytes\n", file.Name, file.isDir(), file.atLocal(), file.Size)
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

type LocalFile struct {
	Path     string
	FileInfo os.FileInfo
}

func (localFile LocalFile) String() string {
	return fmt.Sprintf("%v, %v", localFile.Path, localFile.FileInfo.Name())
}
