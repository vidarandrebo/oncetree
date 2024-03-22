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
	var mut sync.Mutex
	for _, id := range nodeIDs {
		id := id
		wg.Add(1)
		newNode := NewNode(id, nodeMap[id])
		mut.Lock()
		nodes[id] = newNode
		mut.Unlock()
		go func() {
			mut.Lock()
			node := nodes[id]
			mut.Unlock()
			node.Run("")
			wg.Done()
		}()
		mut.Lock()
		go nodes[id].SetNeighboursFromNodeMap(nodeIDs, nodeMap)
		mut.Unlock()
	}
	return nodes, &wg
}

func StartTestNode() (*Node, *sync.WaitGroup) {
	var wg sync.WaitGroup
	id := "10"
	wg.Add(1)
	node := NewNode(id, ":9090")
	go func() {
		node.Run(":9080")
		wg.Done()
	}()
	return node, &wg
}
