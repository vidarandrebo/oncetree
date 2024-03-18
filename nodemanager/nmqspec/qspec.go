package nmqspec

import (
	nmprotos "github.com/vidarandrebo/oncetree/protos/nodemanager"
)

type QSpec struct {
	NumNodes int
}

func (qs *QSpec) PrepareQF(in *nmprotos.PrepareMessage, replies map[uint32]*nmprotos.PromiseMessage) (*nmprotos.PromiseMessage, bool) {
	panic("not implemented")
	return nil, false
}

func (qs *QSpec) AcceptQF(in *nmprotos.AcceptMessage, replies map[uint32]*nmprotos.LearnMessage) (*nmprotos.LearnMessage, bool) {
	panic("not implemented")
	return nil, false
}
