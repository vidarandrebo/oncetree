package nodemanager

import (
	nmprotos "github.com/vidarandrebo/oncetree/protos/nodemanager"
)

type qspec struct {
	NumNodes int
}

func (qs *qspec) PrepareQF(in *nmprotos.PrepareMessage, replies map[uint32]*nmprotos.PromiseMessage) (*nmprotos.PromiseMessage, bool) {
	panic("not implemented")
	return nil, false
}

func (qs *qspec) AcceptQF(in *nmprotos.AcceptMessage, replies map[uint32]*nmprotos.LearnMessage) (*nmprotos.LearnMessage, bool) {
	panic("not implemented")
	return nil, false
}
