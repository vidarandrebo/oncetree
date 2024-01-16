package main

import (
	"github.com/vidarandrebo/oncetree"
)

func main() {
	node := oncetree.NewNode(":8080")
	node.Run()
}
