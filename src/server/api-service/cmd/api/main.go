package main

import (
	"flag"
	"fmt"
	"github.com/linger1216/go-utils/config"
	"gopkg.in/alecthomas/kingpin.v2"
	"os"

	// This Service
	"github.com/linger1216/jelly-doc/src/server/api-service/svc/server"
)

var (
	configFilename = kingpin.Flag("conf", "yaml config file name").Short('c').
		Default("conf/config.yaml").String()
)

func main() {
	// Update addresses if they have been overwritten by flags
	flag.Parse()

	kingpin.Version("0.1.0")
	kingpin.Parse()

	if _, err := os.Stat(*configFilename); err != nil {
		if os.IsNotExist(err) {
			panic(fmt.Sprintf("%s not found", *configFilename))
		} else {
			panic(err)
		}
	}

	reader := config.NewYamlReader(*configFilename)
	server.Run(reader)
}
