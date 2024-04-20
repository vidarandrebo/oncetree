package nmevents

type TreeRecoveredEvent struct {
	FailedNodeID string
}

func NewTreeRecoveredEvent(failedNodeID string) TreeRecoveredEvent {
	return TreeRecoveredEvent{FailedNodeID: failedNodeID}
}
