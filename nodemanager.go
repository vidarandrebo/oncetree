package oncetree

import (
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

func NewPaxos() *NodeManager {
	return &NodeManager{
		paxosData: make(map[int64]*PaxosData),
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
	nm.paxosData[nm.epoch].PrepareMessages[msg.NodeID] = msg
}

func (n *Node) Prepare(ctx gorums.ServerCtx, request *protos.PrepareMessage) {
	n.paxos.StorePrepare(request)
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
