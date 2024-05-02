package nodemanager

import (
	"log/slog"

	"github.com/relab/gorums"
	"github.com/vidarandrebo/oncetree/consts"
	"github.com/vidarandrebo/oncetree/nodemanager/nmenums"
	nmprotos "github.com/vidarandrebo/oncetree/protos/nodemanager"
)

// Join RPC will either accept the node as one of is children or return the address of another node.
// If max fanout is NOT reached, the node will respond OK and include its ID.
// If max fanout is reached, the node will respond NOT OK and include the next address
// of one of its children. The node will alternate between its children in repeated calls to Join.
func (nm *NodeManager) Join(ctx gorums.ServerCtx, request *nmprotos.JoinRequest) (*nmprotos.JoinResponse, error) {
	nm.logger.Info(
		"RPC Join",
		slog.String("id", request.GetNodeID()),
		slog.String("address", request.GetAddress()),
	)
	nm.joinMut.Lock()
	defer nm.joinMut.Unlock()
	response := &nmprotos.JoinResponse{
		OK:          false,
		NodeID:      nm.id,
		Address:     nm.address,
		NextAddress: "",
	}
	children := nm.Children()

	// respond with one of the children's addresses if max fanout has been reached
	if len(children) == consts.Fanout {
		nextPathID := nm.NextJoinID()
		nextPath, _ := nm.neighbours.Get(nextPathID)
		response.OK = false
		response.NextAddress = nextPath.Address
		return response, nil
	}
	response.OK = true
	response.NodeID = nm.id
	nm.AddNeighbour(request.NodeID, request.Address, nmenums.Child)
	return response, nil
}
