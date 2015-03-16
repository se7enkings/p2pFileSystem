package gui
import (
    "fmt"
    "strings"
)

type Node struct {
    Name string
    IsDir bool
    Size int64
    FileHash [32]byte
//    IsAvailable bool
    Children map[string]Node
}

//type FileList struct {
//    List map[string]Node
//}

func (node Node)String() string{
    return node2str(node, 0)
}
func node2str(node Node, space int) string{
    str := ""
    str += strings.Repeat(" ", space) + fmt.Sprintf("%s, %v \n", node.Name, node.IsDir)
    for _, file := range node.Children {
        if file.IsDir {
            str += node2str(file, space + 1)
        }else{
            str += strings.Repeat(" ", space+1) + fmt.Sprintf("%s, %v \n", file.Name, node.IsDir)
        }
    }
    return str
}
