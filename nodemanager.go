package oncetree

import (
	"fmt"
)

type NodeManager struct {
	neighbours   *ConcurrentMap[string, *Neighbour]
	epoch        int64
	nextGorumsID *RWMutex[uint32]
}

func NewNodeManager() *NodeManager {
	return &NodeManager{
		neighbours:   NewConcurrentMap[string, *Neighbour](),
		nextGorumsID: NewRWMutex[uint32](0),
	}
}

func (nm *NodeManager) HandleFailure(nodeID string) {
}

func (nm *NodeManager) GetNeighbours() []KeyValuePair[string, *Neighbour] {
	return nm.neighbours.Entries()
}

func (nm *NodeManager) GetGorumsNeighbourMap() map[string]uint32 {
	IDs := make(map[string]uint32)
	for _, neighbour := range nm.neighbours.Values() {
		IDs[neighbour.Address] = neighbour.GorumsID
	}
	return IDs
}

func (nm *NodeManager) AddNeighbour(nodeID string, address string, role NodeRole) {
	nextID := nm.nextGorumsID.Lock()
	gorumsID := *nextID
	*nextID += 1
	nm.nextGorumsID.Unlock(&nextID)
	neighbour := NewNeighbour(gorumsID, address, role)
	nm.setNeighbour(nodeID, neighbour)
}

func (nm *NodeManager) setNeighbour(nodeID string, neighbour *Neighbour) {
	nm.neighbours.Set(nodeID, neighbour)
}

func (nm *NodeManager) GetNeighbour(nodeID string) (*Neighbour, bool) {
	value, exists := nm.neighbours.Get(nodeID)
	return value, exists
}

func (nm *NodeManager) AllNeighbourIDs() []string {
	return nm.neighbours.Keys()
}

func (nm *NodeManager) AllNeighbourAddrs() []string {
	addresses := make([]string, 0)
	for _, neighbour := range nm.neighbours.Values() {
		addresses = append(addresses, neighbour.Address)
	}
	return addresses
}

func (nm *NodeManager) resolveNodeIDFromAddress(address string) (string, error) {
	for _, neighbour := range nm.neighbours.Entries() {
		if neighbour.Value.Address == address {
			return neighbour.Key, nil
		}
	}
	return "", fmt.Errorf("node with address %s not found", address)
}

func (nm *NodeManager) GetChildren() []*Neighbour {
	children := make([]*Neighbour, 0)
	for _, neighbour := range nm.neighbours.Values() {
		if neighbour.Role == Child {
			children = append(children, neighbour)
		}
	}
	return children
}

func (nm *NodeManager) GetParent() *Neighbour {
	for _, neighbour := range nm.neighbours.Values() {
		if neighbour.Role == Parent {
			return neighbour
		}
	}
	return nil
}
