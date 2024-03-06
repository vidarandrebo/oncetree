package nodemanager

type NeighbourAddedEvent struct {
	NodeID string
}

func NewNeigbourAddedEvent(nodeID string) NeighbourAddedEvent {
	return NeighbourAddedEvent{NodeID: nodeID}
}

type NeighbourRemovedEvent struct {
	NodeID string
}

func NewNeigbourRemovedEvent(nodeID string) NeighbourRemovedEvent {
	return NeighbourRemovedEvent{NodeID: nodeID}
}
