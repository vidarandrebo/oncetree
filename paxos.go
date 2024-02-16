package oncetree

import (
	"sync"

	"github.com/relab/gorums"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/vidarandrebo/oncetree/protos"
)

type Paxos struct {
	data  map[int64]*PaxosData // epoch -> data
	epoch int64
	rwMut sync.RWMutex
}

func NewPaxos() *Paxos {
	return &Paxos{
		data: make(map[int64]*PaxosData),
	}
}

type PaxosData struct {
	// nodeID -> msg
	PrepareMessages map[string]*protos.PrepareMessage
	PromiseMessages map[string]*protos.PromiseMessage
	AcceptMessages  map[string]*protos.AcceptMessage
	LearnMessages   map[string]*protos.LearnMessage
}

func NewPaxosData() *PaxosData {
	return &PaxosData{
		PrepareMessages: make(map[string]*protos.PrepareMessage),
		PromiseMessages: make(map[string]*protos.PromiseMessage),
		AcceptMessages:  make(map[string]*protos.AcceptMessage),
		LearnMessages:   make(map[string]*protos.LearnMessage),
	}
}

func (p *Paxos) StorePrepare(msg *protos.PrepareMessage) {
	p.rwMut.Lock()
	defer p.rwMut.Unlock()
	if _, ok := p.data[p.epoch]; !ok {
		p.data[p.epoch] = NewPaxosData()
	}
	p.data[p.epoch].PrepareMessages[msg.NodeID] = msg
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
