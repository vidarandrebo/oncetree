package failuredetector

type NodeFailedEvent struct {
	NodeID string
}

func NewNodeFailedEvent(nodeID string) NodeFailedEvent {
	return NodeFailedEvent{NodeID: nodeID}
}
