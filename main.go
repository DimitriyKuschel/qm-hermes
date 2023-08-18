package main

import (
	"flag"
	"fmt"
	"os"
	"queue-manager/internal/di"
	"queue-manager/internal/structures"
)

const (
	AppName      = "QueueManager"
	APPNameLower = "queue_manager"
	Version      = "0.0.1"
)

func main() {

	// cli flags
	cf := new(structures.CliFlags)
	flag.StringVar(&cf.ConfigPath, "config", "/etc/"+APPNameLower+"/"+APPNameLower+".yml", "path to config")
	flag.BoolVar(&cf.DebugMode, "debug", false, "debug mode")
	flag.BoolVar(&cf.VersionPrint, "version", false, "print version")
	flag.BoolVar(&cf.Help, "help", false, "show flags")
	flag.BoolVar(&cf.TestMode, "test", false, "test mode")
	flag.Parse()

	if cf.VersionPrint {
		ExitWithCode(fmt.Sprintf("%s Version: %s\n", AppName, Version), 0)
	}

	if cf.Help {
		flag.PrintDefaults()
		os.Exit(0)
	}

	_, err := di.InitApp(cf)
	if err != nil {
		ExitWithCode(fmt.Sprintf("Error app init: %v\n", err), 1)
	}
}

func ExitWithCode(msg string, code int) {
	fmt.Println(msg)
	os.Exit(code)
}
