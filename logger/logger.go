package logger

import (
	"fmt"
	"os"
	"runtime"
)

var debugMode bool = true

func Error(e interface{}) {
	if e != nil {
		fmt.Print("Error: from ")
		pc, _, _, ok := runtime.Caller(1)
		if ok {
			fmt.Print(runtime.FuncForPC(pc).Name(), ": ")
		}
		fmt.Println(e)
		os.Exit(1)
	}
}
func Warning(w interface{}) {
	if w != nil && debugMode {
		fmt.Print("Warning: ")
		pc, _, _, ok := runtime.Caller(1)
		if ok {
			fmt.Print(runtime.FuncForPC(pc).Name(), ": ")
		}
		fmt.Println(w)
	}
}
func Info(i interface{}) {
	if i != nil && debugMode {
		fmt.Printf("Info: ")
		pc, _, _, ok := runtime.Caller(1)
		if ok {
			fmt.Print(runtime.FuncForPC(pc).Name(), ": ")
		}
		fmt.Println(i)
	}
}
