package main

import (
	"github.com/hetianyi/gox/logger"
	"github.com/hetianyi/mysql-2-java-pojo/command"
	"os"
	"runtime"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	args := os.Args[1:]
	for _, arg := range args {
		if arg == "-v" || arg == "--version" {
			args = []string{"--version"}
		}
	}
	newArgs := append([]string{os.Args[0]}, args...)
	// initialize logger
	logConfig := &logger.Config{
		Level:              logger.InfoLevel,
		Write2File:         false,
		AlwaysWriteConsole: true,
	}
	logger.Init(logConfig)

	command.Parse(newArgs)
}
