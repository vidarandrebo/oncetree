package nodemanager

import (
	"log/slog"
	"slices"

	"github.com/vidarandrebo/oncetree/common/hashset"

	"github.com/vidarandrebo/oncetree/storage/sevents"

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
		nm.logger.Error("more than 1 common failure, unrecoverable",
			slog.String("id", request.GetGroupID()))
		panic("more than 1 common failure, unrecoverable")
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

	// do not send data from failed node to new parent
	gossipExcludes := make(map[string]hashset.HashSet[string])
	gossipExcludes[nm.recoveryProcess.newParent] = hashset.New[string]()
	gossipExcludes[nm.recoveryProcess.newParent].Add(nm.recoveryProcess.groupID)

	// do not send data from new parent to existing children
	for _, neighbourID := range nm.NeighbourIDs() {
		if neighbourID != nm.recoveryProcess.groupID {
			gossipExcludes[neighbourID] = hashset.New[string]()
			gossipExcludes[neighbourID].Add(nm.recoveryProcess.newParent)
		}
	}

	nm.eventBus.Execute(sevents.ExcludeFromStorageEvent{
		ReadExcludeIDs: []string{nm.recoveryProcess.newParent}, // we do not read from new parent's data until it is joined with the failed node
		GossipExcludes: gossipExcludes,
	})
	return &nmprotos.LearnMessage{OK: true}, nil
}
