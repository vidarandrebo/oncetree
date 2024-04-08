package main

import (
	"github.com/vidarandrebo/oncetree"
)

var gorumsNodeMap = map[string]uint32{
	":9080": 0,
	":9081": 1,
	":9082": 2,
	":9083": 3,
	":9084": 4,
	":9085": 5,
	":9086": 6,
	":9087": 7,
	":9088": 8,
	":9089": 9,
}

func main() {
	client := oncetree.NewClient(gorumsNodeMap)
	client.Run()
}
