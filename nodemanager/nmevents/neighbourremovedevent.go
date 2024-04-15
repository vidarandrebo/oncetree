package nmevents

type NeighbourRemovedEvent struct {
	NodeID string
}

func NewNeigbourRemovedEvent(nodeID string) NeighbourRemovedEvent {
	return NeighbourRemovedEvent{NodeID: nodeID}
}
