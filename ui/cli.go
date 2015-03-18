package ui

import (
	"fmt"
	"github.com/CRVV/p2pFileSystem/filesystem"
	"github.com/CRVV/p2pFileSystem/settings"
)

func StartCLI() {
	current := filesystem.FileList
	command := ""

Loop:
	for {
		settingValue := settings.GetSettings()
		fmt.Printf("%s@%s %s -->", settingValue.GetUsername(), settingValue.GetGroupName(), current.Name)
		fmt.Scan(&command)
		switch command {
		case "ls":
			fmt.Println(filesystem.Node2str(&current, 0, false))
		case "lstree":
			fmt.Println(filesystem.FileList)
		case "cd":

		case "get":
			filename := ""
			fmt.Scan(&filename)
			fmt.Println(filename)

		case "rm":

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
    help        Print this help
    get FILE    Download file FILE. The FILE must be in current directory
    rm FILE     Remove file FILE. The FILE must be in current directory
    exit        Exit this program

`
