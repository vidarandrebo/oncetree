package nmevents

type NeighbourAddedEvent struct {
	NodeID string
}

func NewNeigbourAddedEvent(nodeID string) NeighbourAddedEvent {
	return NeighbourAddedEvent{NodeID: nodeID}
}
