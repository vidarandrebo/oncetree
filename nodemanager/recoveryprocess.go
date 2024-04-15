package nodemanager

type RecoveryProcess struct {
	GroupID   string
	LeaderID  string
	IsLeader  bool
	NewParent string
}

func NewRecoveryProcess(groupID string, isLeader bool) *RecoveryProcess {
	return &RecoveryProcess{
		GroupID:   groupID,
		LeaderID:  "",
		IsLeader:  isLeader,
		NewParent: "",
	}
}
