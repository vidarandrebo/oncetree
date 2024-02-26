package oncetree_test

import (
	"sync"

	"github.com/vidarandrebo/oncetree"
)

var (
	nodeIDs   = []string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9"}
	nodeAddrs = []string{
		":9080",
		":9081",
		":9082",
		":9083",
		":9084",
		":9085",
		":9086",
		":9087",
		":9088",
		":9089",
	}
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
