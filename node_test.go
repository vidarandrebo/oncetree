package oncetree

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	nodeIDs = []string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9"}
	nodeMap = map[string]string{
		"0": ":0",
		"1": ":1",
		"2": ":2",
		"3": ":3",
		"4": ":4",
		"5": ":5",
		"6": ":6",
		"7": ":7",
		"8": ":8",
		"9": ":9",
	}
)

func TestNode_SetNeighboursFromNodeList_Parent(t *testing.T) {
	node := NewNode("3", ":3")
	node.SetNeighboursFromNodeMap(nodeIDs, nodeMap)
	assert.Equal(t, &Neighbour{Address: ":1", Group: make(map[string]*GroupMember), Role: Parent}, node.nodeManager.GetParent())

	node = NewNode("5", ":5")
	node.SetNeighboursFromNodeMap(nodeIDs, nodeMap)
	assert.Equal(t, &Neighbour{Address: ":2", Group: make(map[string]*GroupMember), Role: Parent}, node.nodeManager.GetParent())
}

func TestNode_SetNeighboursFromNodeList_NoParent(t *testing.T) {
	node := NewNode("0", ":0")
	node.SetNeighboursFromNodeMap(nodeIDs, nodeMap)
	assert.True(t, node.isRoot())
}

func TestNode_SetNeighboursFromNodeList_Children(t *testing.T) {
	node := NewNode("3", ":3")
	node.SetNeighboursFromNodeMap(nodeIDs, nodeMap)

	foundChildren := node.nodeManager.GetChildren()
	assert.Contains(t, foundChildren, &Neighbour{Address: ":7", Group: make(map[string]*GroupMember), Role: Child})
	assert.Contains(t, foundChildren, &Neighbour{Address: ":8", Group: make(map[string]*GroupMember), Role: Child})
}

func TestNode_SetNeighboursFromNodeList_NoChildren(t *testing.T) {
	node := NewNode("5", ":5")
	node.SetNeighboursFromNodeMap(nodeIDs, nodeMap)
	assert.Equal(t, 0, len(node.nodeManager.GetChildren()))
}
