package main

import (
	"github.com/vidarandrebo/once-tree/pkg/node"
)

func main() {
	node := node.NewNode(":8080", ":8081")
	node.Run()
}
