package main

import (
	"sync"

	"github.com/vidarandrebo/oncetree"
)

var (
	nodeIDs = []string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9"}
	nodeMap = map[string]string{
		"0": ":8080",
		"1": ":8081",
		"2": ":8082",
		"3": ":8083",
		"4": ":8084",
		"5": ":8085",
		"6": ":8086",
		"7": ":8087",
		"8": ":8088",
		"9": ":8089",
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
			node.TryUpdateGorumsConfig()
			node.Run()
			wg.Done()
		}()
	}
	wg.Wait()
}
