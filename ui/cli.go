package ui

import (
	"fmt"
	"github.com/CRVV/p2pFileSystem/filesystem"
	"github.com/CRVV/p2pFileSystem/settings"
)

func StartCLI() {
	currentDir := &filesystem.FileList
	paths := make([]string, 0)

Loop:
	for {
		command := ""
		settingValue := settings.GetSettings()
		path := ""
		for _, dir := range paths {
			path += "/"
			path += dir
		}
		if path == "" {
			path = "/"
		}
		fmt.Printf("%s@%s %s --> ", settingValue.GetUsername(), settingValue.GetGroupName(), path)
		fmt.Scan(&command)
		switch command {
		case "ls":
			fmt.Println(filesystem.Node2str(currentDir, 0, false))
		case "lstree":
			filesystem.FlMutex.Lock()
			fmt.Println(filesystem.FileList)
			filesystem.FlMutex.Unlock()
		case "cd":
			name := ""
			fmt.Scan(&name)
			newDir, ok := currentDir.Children[name]
			switch {
			case name == ".":
			case name == "..":
				if len(paths) > 0 {
					paths = paths[0 : len(paths)-1]
				}
				currentDir = newDir
			case ok && newDir.IsDir:
				paths = append(paths, name)
				currentDir = newDir
			case ok && !newDir.IsDir:
				fmt.Printf("not a directory: %s\n", name)
			default:
				fmt.Printf("no such file or directory: %s\n", name)
			}
		case "get":
			name := ""
			fmt.Scan(&name)
			file, ok := currentDir.Children[name]
			switch {
			case !ok:
				fmt.Printf("no such file or directory: %s\n", name)
			case file.IsDir:
				fmt.Printf("not a file: %s\n", name)
			case file.AtLocal:
				fmt.Println("download complete")
			case !file.AtLocal:
				filesystem.GetFile(file.FileHash)
				filesystem.Init()
			}
		case "rm":
			name := ""
			fmt.Scan(&name)
			file, ok := currentDir.Children[name]
			switch {
			case !ok:
				fmt.Printf("no such file or directory: %s\n", name)
			case file.IsDir:
				fmt.Printf("not a file: %s\n", name)
			case !file.AtLocal:
				fmt.Printf("not a local file: %s\n", name)
			case !file.IsDir:
				filesystem.RemoveLocalFile(file.FileHash)
				filesystem.Init()
			}
		case "help":
			fmt.Print(HELP)
		case "exit":
			break Loop
		default:
			fmt.Printf("command not found: %s\n", command)
		}
	}
}

const HELP string = `The commands are:
    ls          List files in current directory
    lstree      List the whole directory tree
    cd DIR      change directory to DIR
    get FILE    Download file FILE. The FILE must be in current directory
    rm FILE     Remove file FILE. The FILE must be in current directory
    help        Print this help
    exit        Exit this program
`
