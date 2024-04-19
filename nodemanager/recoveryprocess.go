package nodemanager

import "sync"

type RecoveryProcess struct {
	isActive  bool
	isLeader  bool
	mut       sync.RWMutex
	groupID   string
	leaderID  string
	newParent string
}

func (rp *RecoveryProcess) start(groupID string) {
	rp.isActive = true
	rp.groupID = groupID
}

func (rp *RecoveryProcess) stop() {
	rp.isActive = false
	rp.isLeader = false
	rp.groupID = ""
	rp.leaderID = ""
	rp.newParent = ""
}
