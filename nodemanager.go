package oncetree

import (
	"fmt"
	"sync"
)

type NodeManager struct {
	neighbours map[string]*Neighbour
	epoch      int64
	rwMut      sync.RWMutex
}

func NewNodeManager() *NodeManager {
	return &NodeManager{
		neighbours: make(map[string]*Neighbour),
	}
}

func (nm *NodeManager) HandleFailure(nodeID string) {
}

func (nm *NodeManager) GetNeighbours() map[string]*Neighbour {
	return nm.neighbours
}

func (nm *NodeManager) SetNeighbour(nodeID string, neighbour *Neighbour) {
	nm.neighbours[nodeID] = neighbour
}

func (nm *NodeManager) GetNeighbour(nodeID string) (*Neighbour, bool) {
	value, exists := nm.neighbours[nodeID]
	return value, exists
}

func (nm *NodeManager) AllNeighbourIDs() []string {
	IDs := make([]string, 0)
	for id := range nm.neighbours {
		IDs = append(IDs, id)
	}
	return IDs
}

func (nm *NodeManager) AllNeighbourAddrs() []string {
	addresses := make([]string, 0)
	for _, neighbour := range nm.neighbours {
		addresses = append(addresses, neighbour.Address)
	}
	return addresses
}

func (nm *NodeManager) resolveNodeIDFromAddress(address string) (string, error) {
	for id, neighbour := range nm.neighbours {
		if neighbour.Address == address {
			return id, nil
		}
	}
	return "", fmt.Errorf("node with address %s not found", address)
}

func (nm *NodeManager) GetChildren() []*Neighbour {
	children := make([]*Neighbour, 0)
	for _, neighbour := range nm.neighbours {
		if neighbour.Role == Child {
			children = append(children, neighbour)
		}
	}
	return children
}

func (nm *NodeManager) GetParent() *Neighbour {
	for _, neighbour := range nm.neighbours {
		if neighbour.Role == Parent {
			return neighbour
		}
	}
	return nil
}
