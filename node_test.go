package oncetree_test

import (
	"sync"

	"github.com/vidarandrebo/oncetree"
)

var (
	nodeIDs   = []string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9"}
	nodeAddrs = []string{
		":8080",
		":8081",
		":8082",
		":8083",
		":8084",
		":8085",
		":8086",
		":8087",
		":8088",
		":8089",
	}
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

//func TestNode_SetNeighboursFromNodeList_Parent(t *testing.T) {
//	node := oncetree.NewNode("3", ":3")
//	node.SetNeighboursFromNodeMap(nodeIDs, nodeMap)
//	assert.Equal(t, &oncetree.Neighbour{Address: ":1", Group: make(map[string]*oncetree.GroupMember), Role: oncetree.Parent}, node.nodeManager.GetParent())
//
//	node = NewNode("5", ":5")
//	node.SetNeighboursFromNodeMap(nodeIDs, nodeMap)
//	assert.Equal(t, &Neighbour{Address: ":2", Group: make(map[string]*GroupMember), Role: Parent}, node.nodeManager.GetParent())
//}
//
//func TestNode_SetNeighboursFromNodeList_NoParent(t *testing.T) {
//	node := NewNode("0", ":0")
//	node.SetNeighboursFromNodeMap(nodeIDs, nodeMap)
//	assert.True(t, node.isRoot())
//}
//
//func TestNode_SetNeighboursFromNodeList_Children(t *testing.T) {
//	node := NewNode("3", ":3")
//	node.SetNeighboursFromNodeMap(nodeIDs, nodeMap)
//
//	foundChildren := node.nodeManager.GetChildren()
//	assert.Contains(t, foundChildren, &Neighbour{Address: ":7", Group: make(map[string]*GroupMember), Role: Child})
//	assert.Contains(t, foundChildren, &Neighbour{Address: ":8", Group: make(map[string]*GroupMember), Role: Child})
//}
//
//func TestNode_SetNeighboursFromNodeList_NoChildren(t *testing.T) {
//	node := NewNode("5", ":5")
//	node.SetNeighboursFromNodeMap(nodeIDs, nodeMap)
//	assert.Equal(t, 0, len(node.nodeManager.GetChildren()))
//}

func StartTestNodes() (map[string]*oncetree.Node, *sync.WaitGroup) {
	var wg sync.WaitGroup
	nodes := make(map[string]*oncetree.Node)
	for _, id := range nodeIDs {
		id := id
		wg.Add(1)
		go func() {
			node := oncetree.NewNode(id, nodeMap[id])
			node.SetNeighboursFromNodeMap(nodeIDs, nodeMap)
			nodes[id] = node
			node.Run()
			wg.Done()
		}()
	}
	return nodes, &wg
}
