package nodemanager

import (
	"log/slog"
	"slices"
	"sort"

	"github.com/relab/gorums"
	"github.com/vidarandrebo/oncetree/nodemanager/nmenums"
	nmprotos "github.com/vidarandrebo/oncetree/protos/nodemanager"
)

func (nm *NodeManager) Prepare(ctx gorums.ServerCtx, request *nmprotos.PrepareMessage) (*nmprotos.PromiseMessage, error) {
	nm.logger.Info("Prepare RPC",
		"id", request.NodeID)
	node, ok := nm.Neighbour(request.GetGroupID())
	if !ok {
		// this could mean a delayed prepare message from one of the nodes.
		return &nmprotos.PromiseMessage{OK: false}, nil
	}

	if node.Group.epoch != request.Epoch {
		panic("node's group info is not up to date, unrecoverable")
	}
	nm.recoveryProcess.mut.Lock()
	defer nm.recoveryProcess.mut.Unlock()

	if (nm.recoveryProcess.isActive) && (nm.recoveryProcess.groupID != request.GetGroupID()) {
		nm.logger.Error("more than 1 concurrent failure, unrecoverable",
			slog.String("id", request.GetGroupID()))
		panic("more than 1 concurrent failure, unrecoverable")
	}
	if !nm.recoveryProcess.isActive {
		nm.recoveryProcess.start(request.GetGroupID())
	}

	for _, member := range node.Group.members {
		if member.Role == nmenums.Parent {
			if member.ID == request.GetNodeID() {
				nm.logger.Info("member is Parent, send leader promise",
					slog.String("id", request.GetNodeID()))
				nm.recoveryProcess.leaderID = request.GetNodeID()

				// leader is added as a recovery node here
				// it will be converted into a parent node in the end of the recovery session
				_, neighbourExists := nm.Neighbour(member.ID)
				if !neighbourExists {
					nm.AddNeighbour(member.ID, member.Address, nmenums.Recovery)
				}
				return &nmprotos.PromiseMessage{OK: true}, nil
			} else {
				nm.logger.Info("member is not Parent, deny as leader",
					slog.String("id", request.GetNodeID()))
				return &nmprotos.PromiseMessage{OK: false}, nil
			}
		}
	}

	// sort member IDs alphabetically, then return OK if the id from the request is the first in the collection
	nm.logger.Info("failed node has no parent, using lowest ID node as leader")
	memberIDs := node.GroupMemberIDs()
	sort.Strings(memberIDs)
	rank := slices.Index(memberIDs, request.GetNodeID())
	if rank == 0 {
		nm.logger.Info("member has lowest ID, send leader promise",
			slog.String("id", request.GetNodeID()))
		nm.recoveryProcess.leaderID = request.GetNodeID()

		// find member from request in collection
		for _, member := range node.Group.members {
			if member.ID != request.GetNodeID() {
				continue
			}
			// leader is added as a recovery node here
			// it will be converted into a parent node in the end of the recovery session
			_, neighbourExists := nm.Neighbour(member.ID)
			if !neighbourExists {
				nm.AddNeighbour(member.ID, member.Address, nmenums.Recovery)
			}
		}
		return &nmprotos.PromiseMessage{OK: true}, nil
	}
	return &nmprotos.PromiseMessage{OK: false}, nil
}
