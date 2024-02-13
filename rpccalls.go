package oncetree

import (
	"github.com/relab/gorums"
	"github.com/vidarandrebo/oncetree/protos"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (n *Node) Write(ctx gorums.ServerCtx, request *protos.WriteRequest) (response *emptypb.Empty, err error) {
	n.mut.Lock()
	n.timestamp++
	n.keyValueStorage.WriteValue(n.id, request.GetKey(), request.GetValue(), n.timestamp)
	n.mut.Unlock()
	go func() {
		n.sendGossip(n.id, request.GetKey())
	}()
	return &emptypb.Empty{}, nil
}

func (n *Node) Read(ctx gorums.ServerCtx, request *protos.ReadRequest) (response *protos.ReadResponse, err error) {
	value, err := n.keyValueStorage.ReadValue(request.Key)
	if err != nil {
		return &protos.ReadResponse{Value: 0}, err
	}
	return &protos.ReadResponse{Value: value}, nil
}

func (n *Node) ReadAll(ctx gorums.ServerCtx, request *protos.ReadRequest) (response *protos.ReadAllResponse, err error) {
	value, err := n.keyValueStorage.ReadValue(request.Key)
	if err != nil {
		return &protos.ReadAllResponse{Value: nil}, err
	}
	return &protos.ReadAllResponse{Value: map[string]int64{n.id: value}}, nil
}

func (n *Node) SetGroupMember(ctx gorums.ServerCtx, request *protos.GroupInfo) {
	n.logger.Printf("Adding node %s to group %s", request.NeighbourID, request.NodeID)
	n.neighbours[request.NodeID].Group[request.NeighbourID] = request.NeighbourAddress
}

func (n *Node) PrintState(ctx gorums.ServerCtx, request *emptypb.Empty) (response *emptypb.Empty, err error) {
	n.logger.Println(n.keyValueStorage)
	return &emptypb.Empty{}, nil
}

func (n *Node) Gossip(ctx gorums.ServerCtx, request *protos.GossipMessage) (response *emptypb.Empty, err error) {
	n.logger.Printf("received gossip %v", request)
	n.mut.Lock()
	updated := n.keyValueStorage.WriteValue(request.GetNodeID(), request.GetKey(), request.GetValue(), request.GetTimestamp())
	n.mut.Unlock()
	if updated {
		go func() {
			n.sendGossip(request.NodeID, request.GetKey())
		}()
	}
	return &emptypb.Empty{}, nil
}

func (n *Node) Heartbeat(ctx gorums.ServerCtx, request *protos.HeartbeatMessage) {
	// n.logger.Printf("received heartbeat from %s", request.GetNodeID())
	n.failureDetector.RegisterHeartbeat(request.GetNodeID())
}
