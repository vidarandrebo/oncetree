package nodemanager

import (
	"github.com/relab/gorums"
	"github.com/vidarandrebo/oncetree/nodemanager/nmenums"
	"github.com/vidarandrebo/oncetree/nodemanager/nmevents"
	nmprotos "github.com/vidarandrebo/oncetree/protos/nodemanager"
)

func (nm *NodeManager) Commit(ctx gorums.ServerCtx, request *nmprotos.CommitMessage) {
	nm.logger.Info("Commit RPC",
		"failedID", request.GetGroupID())
	nm.recoveryProcess.mut.Lock()
	defer nm.recoveryProcess.mut.Unlock()
	if request.GetGroupID() != nm.recoveryProcess.groupID {
		panic("request's group is not same as stored in recoveryprocess")
	}
	newParentID := nm.recoveryProcess.newParent
	newParent, ok := nm.neighbours.Get(newParentID)
	if !ok {
		panic("could not find new parent in neighbour map")
	}
	newParent.Role = nmenums.Parent
	nm.neighbours.Delete(request.GetGroupID())
	nm.blackList.Add(request.GetGroupID())
	nm.eventBus.PushEvent(nmevents.NewNeigbourRemovedEvent(request.GetGroupID()))

	newGorumsNeighbourMap := nm.GorumsNeighbourMap()
	nm.gorumsProvider.ResetWithNewNodes(newGorumsNeighbourMap)

	nm.eventBus.PushEvent(nmevents.NewNeighbourAddedEvent(newParentID, "commit", nmenums.Parent))
	nm.recoveryProcess.stop()
	nm.eventBus.PushTask(nm.SendGroupInfo)
}
