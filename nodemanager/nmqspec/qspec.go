package nmqspec

import (
	nmprotos "github.com/vidarandrebo/oncetree/protos/nodemanager"
)

type QSpec struct {
	NumNodes int
}

func (qs *QSpec) PrepareQF(in *nmprotos.PrepareMessage, replies map[uint32]*nmprotos.PromiseMessage) (*nmprotos.PromiseMessage, bool) {
	// all nodes have to answer this prepare
	if len(replies) != qs.NumNodes {
		return &nmprotos.PromiseMessage{OK: false}, false
	}

	for _, reply := range replies {
		if !reply.OK {
			return &nmprotos.PromiseMessage{OK: false}, true
		}
	}
	return &nmprotos.PromiseMessage{OK: true}, true
}

func (qs *QSpec) AcceptQF(in *nmprotos.AcceptMessage, replies map[uint32]*nmprotos.LearnMessage) (*nmprotos.LearnMessage, bool) {
	// all nodes have to answer this accept
	if len(replies) != qs.NumNodes {
		return &nmprotos.LearnMessage{OK: false}, false
	}

	for _, reply := range replies {
		if !reply.OK {
			return &nmprotos.LearnMessage{OK: false}, true
		}
	}
	return &nmprotos.LearnMessage{OK: true}, true
}
