package ui

import (
	"bufio"
	"fmt"
	"github.com/CRVV/p2pFileSystem/filesystem"
	"github.com/CRVV/p2pFileSystem/remote"
	"github.com/CRVV/p2pFileSystem/settings"
	"os"
	"strings"
)

func StartCLI() {
	var currentDir *filesystem.Node
	var paths []string
Loop:
	for {
		settingValue := settings.GetSettings()
		path := ""
		fileList := filesystem.GetFileList()
		fileList.RLock()
		currentDir = fileList.N
		for _, dir := range paths {
			var ok bool
			currentDir, ok = currentDir.Children[dir]
			if !ok {
				currentDir = fileList.N
				paths = nil
				path = ""
				break
			}
			path += "/"
			path += dir
		}
		fileList.RUnlock()
		if path == "" {
			path = "/"
		}
		fmt.Printf("%s@%s %s --> ", settingValue.GetUsername(), settingValue.GetGroupName(), path)

		cmdReader := bufio.NewReaderSize(os.Stdin, 128)
		cmd, _, _ := cmdReader.ReadLine()
		command := strings.Split(string(cmd), " ")
		name := ""
		if len(command) > 1 {
			name = command[1]
		}
		switch command[0] {
		case "ls":
			fmt.Println(filesystem.Node2str(currentDir, 0, false))
		case "lstree":
			fileList.RLock()
			fmt.Println(fileList.N)
			fileList.RUnlock()
		case "cd":
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
			file, ok := currentDir.Children[name]
			switch {
			case !ok:
				fmt.Printf("no such file or directory: %s\n", name)
			case file.IsDir:
				fmt.Printf("not a file: %s\n", name)
			case file.AtLocal:
				fmt.Println("download complete")
			case !file.AtLocal:
				err := remote.DownloadFile(file.FileHash)
				if err != nil {
					fmt.Println(err)
				} else {
					fmt.Println("Download complete")
				}
			}
		case "rm":
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
				fmt.Println("remove file complete")
			}
		case "help":
			fmt.Print(HELP)
		case "exit":
			break Loop
		default:
			if command[0] != "" {
				fmt.Printf("command not found: %s\n", command[0])
			}
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
