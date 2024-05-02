package oncetree

import (
	"io"
	"sync"
	"time"

	"github.com/vidarandrebo/oncetree/consts"
)

func StartTestNodes(discardLogs bool) (map[string]*Node, *sync.WaitGroup) {
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
	var readyWg sync.WaitGroup
	var wg sync.WaitGroup
	nodes := make(map[string]*Node)
	var mut sync.Mutex
	for _, id := range nodeIDs {
		readyWg.Add(1)
		id := id
		wg.Add(1)
		newNode := NewNode(id, nodeMap[id], io.Discard)
		mut.Lock()
		nodes[id] = newNode
		mut.Unlock()
		go func() {
			mut.Lock()
			node := nodes[id]
			mut.Unlock()
			node.Run(GetParentFromNodeMap(id, nodeIDs, nodeMap), readyWg.Done)
			wg.Done()
		}()
	}
	readyWg.Wait()
	time.Sleep(consts.StartupDelay)
	return nodes, &wg
}

func StartTestNode(discardLogs bool) (*Node, *sync.WaitGroup) {
	var readyWg sync.WaitGroup
	readyWg.Add(1)
	var wg sync.WaitGroup
	id := "10"
	wg.Add(1)
	node := NewNode(id, ":9090", io.Discard)
	go func() {
		node.Run(":9080", readyWg.Done)
		wg.Done()
	}()
	time.Sleep(consts.StartupDelay)
	return node, &wg
}

// GetParentFromNodeMap assumes a binary tree as slice where a nodes children are at index 2i+1 and 2i+2
// Uses both a slice and a maps to ensure consistent iteration order
func GetParentFromNodeMap(id string, nodeIDs []string, nodes map[string]string) string {
	for i, nodeID := range nodeIDs {
		// find n as a child of current node -> current node is n's parent
		if len(nodeIDs) > (2*i+1) && nodeIDs[2*i+1] == id {
			return nodes[nodeID]
		}
		if len(nodeIDs) > (2*i+2) && nodeIDs[2*i+2] == id {
			return nodes[nodeID]
		}
	}
	return ""
}
