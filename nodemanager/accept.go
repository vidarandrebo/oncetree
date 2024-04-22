package nodemanager

import (
	"log/slog"
	"slices"

	"github.com/relab/gorums"
	nmprotos "github.com/vidarandrebo/oncetree/protos/nodemanager"
)

func (nm *NodeManager) Accept(ctx gorums.ServerCtx, request *nmprotos.AcceptMessage) (*nmprotos.LearnMessage, error) {
	nm.logger.Info("nm.Accept RPC",
		"id", request.GetNodeID())
	nm.recoveryProcess.mut.Lock()
	defer nm.recoveryProcess.mut.Unlock()

	// no ongoing recovery process
	if !nm.recoveryProcess.isActive {
		return &nmprotos.LearnMessage{OK: false}, nil
	}
	// requester is not leader
	if nm.recoveryProcess.leaderID != request.GetNodeID() {
		return &nmprotos.LearnMessage{OK: false}, nil
	}

	if nm.recoveryProcess.groupID != request.GetGroupID() {
		nm.logger.Error("more than 1 concurrent failure, unrecoverable",
			slog.String("id", request.GetGroupID()))
		panic("more than 1 concurrent failure, unrecoverable")
	}

	failedNode, ok := nm.neighbours.Get(request.GetGroupID())
	if !ok {
		return &nmprotos.LearnMessage{OK: false}, nil
	}
	newParent := request.GetNewParent()[nm.id]

	// new parent has to be part of the recovery group
	containsNewParent := slices.Contains(failedNode.GroupMemberIDs(), newParent)
	if !containsNewParent {
		return &nmprotos.LearnMessage{OK: false}, nil
	}

	nm.recoveryProcess.newParent = newParent
	return &nmprotos.LearnMessage{OK: true}, nil
}
