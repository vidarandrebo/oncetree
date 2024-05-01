package nodemanager

import (
	"log/slog"

	"github.com/relab/gorums"
	"github.com/vidarandrebo/oncetree/nodemanager/nmenums"
	nmprotos "github.com/vidarandrebo/oncetree/protos/nodemanager"
)

func (nm *NodeManager) GroupInfo(ctx gorums.ServerCtx, request *nmprotos.GroupInfoMessage) {
	// message arrive one by one from same client in gorums, so should not need to lock for epoch compare
	node, ok := nm.neighbours.Get(request.GetGroupID())

	// Group exists and is up to date
	if ok && node.Group.epoch >= request.GetEpoch() {
		return
	}
	if !ok {
		// this warning will trigger if another node has sent its ready signal to node 'a' before this node has
		// sent its ready signal to 'a' in the join process. This is expected when many nodes are
		// joining the network concurrently. All group info will be distributed regardless of this.
		nm.logger.Warn("received group info from unknown node",
			slog.String("id", request.GetGroupID()))
		return
	}

	newMembers := make([]GroupMember, 0)
	for _, member := range request.GetMembers() {
		newMembers = append(newMembers, NewGroupMember(
			member.GetID(),
			member.GetAddress(),
			nmenums.NodeRole(member.GetRole())),
		)
	}
	node.Group.epoch = request.GetEpoch()
	node.Group.members = newMembers
	nm.logger.Info("group updated",
		slog.String("id", request.GetGroupID()),
		slog.String("group", node.Group.String()))
}
