package oncetree

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

var nodeIDs = []string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9"}
var nodeMap = map[string]string{
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

func TestNode_SetNeighboursFromNodeList_Parent(t *testing.T) {
	node := NewNode("3", ":3")
	node.SetNeighboursFromNodeMap(nodeIDs, nodeMap)
	assert.Equal(t, map[string]string{"1": ":1"}, node.parent)

	node = NewNode("5", ":5")
	node.SetNeighboursFromNodeMap(nodeIDs, nodeMap)
	assert.Equal(t, map[string]string{"2": ":2"}, node.parent)
}
func TestNode_SetNeighboursFromNodeList_NoParent(t *testing.T) {
	node := NewNode("0", ":0")
	node.SetNeighboursFromNodeMap(nodeIDs, nodeMap)
	assert.True(t, node.isRoot())
}
func TestNode_SetNeighboursFromNodeList_Children(t *testing.T) {
	node := NewNode("3", ":3")
	node.SetNeighboursFromNodeMap(nodeIDs, nodeMap)
	assert.Equal(t, map[string]string{"7": ":7", "8": ":8"}, node.children)
}
func TestNode_SetNeighboursFromNodeList_NoChildren(t *testing.T) {
	node := NewNode("5", ":5")
	node.SetNeighboursFromNodeMap(nodeIDs, nodeMap)
	assert.Equal(t, 0, len(node.children))
}
