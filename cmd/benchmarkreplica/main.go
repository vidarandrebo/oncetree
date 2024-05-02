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
	knownAddr := flag.String("knownAddr", "", "IP address of one of the nodes in the network")
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

	// Address is a docker dns trick
	node := oncetree.NewNode(id, fmt.Sprintf("%s:8080", id), file)

	node.Run(*knownAddr, nil)
}
