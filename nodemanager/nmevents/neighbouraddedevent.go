package nmevents

import (
	"github.com/vidarandrebo/oncetree/nodemanager/nmenums"
)

type NeighbourAddedEvent struct {
	NodeID string
	Role   nmenums.NodeRole
}

func NewNeighbourAddedEvent(nodeID string, role nmenums.NodeRole) NeighbourAddedEvent {
	return NeighbourAddedEvent{NodeID: nodeID, Role: role}
}
