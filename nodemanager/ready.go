package nodemanager

import (
	"fmt"
	"log/slog"

	"github.com/relab/gorums"
	"github.com/vidarandrebo/oncetree/nodemanager/nmevents"
	nmprotos "github.com/vidarandrebo/oncetree/protos/nodemanager"
)

// Ready is used to signal to the parent that the newly joined node has created a gorums config and is ready to participate in the protocol.
// A NeighbourReadyEvent is pushed to the eventbus
func (nm *NodeManager) Ready(ctx gorums.ServerCtx, request *nmprotos.ReadyMessage) (*nmprotos.ReadyMessage, error) {
	nm.logger.Debug(
		"RPC Ready",
		slog.String("id", request.GetNodeID()),
	)
	// check node exists
	_, ok := nm.Neighbour(request.GetNodeID())
	if !ok {
		return &nmprotos.ReadyMessage{OK: false, NodeID: nm.id}, fmt.Errorf("node %s is not joined to this node", request.GetNodeID())
	}

	nm.eventBus.PushEvent(nmevents.NewNeighbourReadyEvent(request.GetNodeID()))
	return &nmprotos.ReadyMessage{OK: true, NodeID: nm.id}, nil
}
