package nmevents

type NeighbourReadyEvent struct {
	NodeID string
}

func NewNeighbourReadyEvent(nodeID string) NeighbourReadyEvent {
	return NeighbourReadyEvent{NodeID: nodeID}
}
