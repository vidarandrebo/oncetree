package oncetree

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

var nodeList = []string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9"}

func TestNode_SetNeighboursFromNodeList_Parent(t *testing.T) {
	node := NewNode("3", ":3")
	node.SetNeighboursFromNodeList(nodeList)
	assert.Equal(t, "1", node.parent)

	node = NewNode("5", ":5")
	node.SetNeighboursFromNodeList(nodeList)
	assert.Equal(t, "2", node.parent)
}
func TestNode_SetNeighboursFromNodeList_NoParent(t *testing.T) {
	node := NewNode("0", ":0")
	node.SetNeighboursFromNodeList(nodeList)
	assert.True(t, node.IsRoot())
}
func TestNode_SetNeighboursFromNodeList_Children(t *testing.T) {
	node := NewNode("3", ":3")
	node.SetNeighboursFromNodeList(nodeList)
	assert.Equal(t, []string{"7", "8"}, node.children)
}
func TestNode_SetNeighboursFromNodeList_NoChildren(t *testing.T) {
	node := NewNode("5", ":5")
	node.SetNeighboursFromNodeList(nodeList)
	assert.Equal(t, 0, len(node.children))
}
