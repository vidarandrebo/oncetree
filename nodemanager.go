package oncetree

import (
	"fmt"
	"sync"

	"github.com/relab/gorums"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/vidarandrebo/oncetree/protos"
)

type NodeManager struct {
	neighbours map[string]*Neighbour
	paxosData  map[int64]*PaxosData // epoch -> paxosData
	epoch      int64
	rwMut      sync.RWMutex
}

func NewNodeManager() *NodeManager {
	return &NodeManager{
		neighbours: make(map[string]*Neighbour),
		paxosData:  make(map[int64]*PaxosData),
	}
}

type PaxosData struct {
	// nodeID -> msg
	PrepareMessages chan *protos.PrepareMessage
	PromiseMessages chan *protos.PromiseMessage
	AcceptMessages  chan *protos.AcceptMessage
	LearnMessages   chan *protos.LearnMessage
}

func NewPaxosData() *PaxosData {
	return &PaxosData{
		PrepareMessages: make(chan *protos.PrepareMessage),
		PromiseMessages: make(chan *protos.PromiseMessage),
		AcceptMessages:  make(chan *protos.AcceptMessage),
		LearnMessages:   make(chan *protos.LearnMessage),
	}
}

func (nm *NodeManager) HandleFailure(nodeID string) {
}

func (nm *NodeManager) GetNeighbours() map[string]*Neighbour {
	return nm.neighbours
}
func (nm *NodeManager) SetNeighbour(nodeID string, neighbour *Neighbour) {
	nm.neighbours[nodeID] = neighbour
}
func (nm *NodeManager) GetNeighbour(nodeID string) (*Neighbour, bool) {
	value, exists := nm.neighbours[nodeID]
	return value, exists
}
func (nm *NodeManager) AllNeighbourIDs() []string {
	IDs := make([]string, 0)
	for id := range nm.neighbours {
		IDs = append(IDs, id)
	}
	return IDs
}
func (nm *NodeManager) AllNeighbourAddrs() []string {
	addresses := make([]string, 0)
	for _, neighbour := range nm.neighbours {
		addresses = append(addresses, neighbour.Address)
	}
	return addresses
}

func (nm *NodeManager) StorePrepare(msg *protos.PrepareMessage) {
	nm.rwMut.Lock()
	defer nm.rwMut.Unlock()
	if _, ok := nm.paxosData[nm.epoch]; !ok {
		nm.paxosData[nm.epoch] = NewPaxosData()
	}
	//nm.paxosData[nm.epoch].PrepareMessages[msg.NodeID] = msg
}
func (nm *NodeManager) resolveNodeIDFromAddress(address string) (string, error) {
	for id, neighbour := range nm.neighbours {
		if neighbour.Address == address {
			return id, nil
		}
	}
	return "", fmt.Errorf("node with address %s not found", address)
}
func (nm *NodeManager) GetChildren() []*Neighbour {
	children := make([]*Neighbour, 0)
	for _, neighbour := range nm.neighbours {
		if neighbour.Role == Child {
			children = append(children, neighbour)
		}
	}
	return children
}

func (nm *NodeManager) GetParent() *Neighbour {
	for _, neighbour := range nm.neighbours {
		if neighbour.Role == Parent {
			return neighbour
		}
	}
	return nil
}
func (nm *NodeManager) SetGroupMember(groupID string, memberID string, groupMember *GroupMember) {
	nm.neighbours[groupID].Group[memberID] = groupMember
}

func (n *Node) Prepare(ctx gorums.ServerCtx, request *protos.PrepareMessage) {
	//n.paxos.StorePrepare(request)
}

func (n *Node) Promise(ctx gorums.ServerCtx, request *protos.PromiseMessage) (response *emptypb.Empty, err error) {
	n.logger.Panicf("not implemented")
	return nil, err
}

func (n *Node) Accept(ctx gorums.ServerCtx, request *protos.AcceptMessage) {
	n.logger.Panicf("not implemented")
}

func (n *Node) Learn(ctx gorums.ServerCtx, request *protos.LearnMessage) {
	n.logger.Panicf("not implemented")
}

func (n *Node) SetGroupMember(ctx gorums.ServerCtx, request *protos.GroupInfo) {
	n.logger.Printf("Adding node %s to group %s", request.NeighbourID, request.NodeID)
	groupMember := NewGroupMember(request.NeighbourAddress, NodeRole(request.NeighbourRole))
	n.nodeManager.SetGroupMember(request.NodeID, request.NeighbourID, groupMember)
}
