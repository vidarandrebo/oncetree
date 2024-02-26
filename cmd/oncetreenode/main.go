package main

import (
	"sync"

	"github.com/vidarandrebo/oncetree"
)

var (
	nodeIDs = []string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9"}
	nodeMap = map[string]string{
		"0": ":9080",
		"1": ":9081",
		"2": ":9082",
		"3": ":9083",
		"4": ":9084",
		"5": ":9085",
		"6": ":9086",
		"7": ":9087",
		"8": ":9088",
		"9": ":9089",
	}
)

func main() {
	var wg sync.WaitGroup
	for _, id := range nodeIDs {
		id := id
		wg.Add(1)
		go func() {
			node := oncetree.NewNode(id, nodeMap[id])
			node.SetNeighboursFromNodeMap(nodeIDs, nodeMap)
			node.Run()
			wg.Done()
		}()
	}
	wg.Wait()
}
