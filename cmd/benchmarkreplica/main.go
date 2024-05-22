package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/vidarandrebo/oncetree"
	"github.com/vidarandrebo/oncetree/consts"
)

func main() {
	runtime.GOMAXPROCS(1)
	knownAddress := flag.String("known-address", "", "IP address of one of the nodes in the network")
	address := flag.String("address", "", "IP address to serve on")
	port := flag.String("port", "", "Port to serve on")
	flag.Parse()

	id, err := os.Hostname()
	if err != nil {
		panic("could not get hostname")
	}
	fileName := filepath.Join(consts.LogFolder, fmt.Sprintf("replica_%s.log", id))
	file, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o666)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// id is passed as the address, only works where hostname lookup is configured
	serveAddress := *address
	if serveAddress == "" {
		serveAddress = id
	}
	servePort := *port
	if servePort == "" {
		servePort = "8080"
	}

	node := oncetree.NewNode(id, serveAddress, servePort, file)

	node.Run(*knownAddress, nil)
}
