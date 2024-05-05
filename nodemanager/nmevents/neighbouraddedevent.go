package nmevents

import (
	"github.com/vidarandrebo/oncetree/nodemanager/nmenums"
)

type NeighbourAddedEvent struct {
	NodeID  string
	Address string
	Role    nmenums.NodeRole
}

func NewNeighbourAddedEvent(nodeID string, address string, role nmenums.NodeRole) NeighbourAddedEvent {
	return NeighbourAddedEvent{NodeID: nodeID, Address: address, Role: role}
}
