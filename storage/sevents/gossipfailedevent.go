package sevents

type GossipFailedEvent struct {
	NodeID string
}

func NewGossipFailedEvent(nodeID string) GossipFailedEvent {
	return GossipFailedEvent{
		NodeID: nodeID,
	}
}
