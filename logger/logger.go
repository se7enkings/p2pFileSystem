package logger

import (
	"fmt"
	"os"
)

var debugMode bool = true

func Error(e interface{}) {
	if e != nil && debugMode {
		fmt.Printf("Error: %s\n", e)
		os.Exit(1)
	}
}
func Warning(w interface{}) {
	if w != nil && debugMode {
		fmt.Printf("Warning: %s\n", w)
	}
}
func Info(i interface{}) {
	if i != nil && debugMode {
		fmt.Printf("Info: %s\n", i)
	}
}
