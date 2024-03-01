package oncetree

import (
	"sync"
)

func StartTestNodes() (map[string]*Node, *sync.WaitGroup) {
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
	var wg sync.WaitGroup
	nodes := make(map[string]*Node)
	for _, id := range nodeIDs {
		id := id
		wg.Add(1)
		go func() {
			node := NewNode(id, nodeMap[id])
			node.SetNeighboursFromNodeMap(nodeIDs, nodeMap)
			nodes[id] = node
			node.Run()
			wg.Done()
		}()
	}
	return nodes, &wg
}
