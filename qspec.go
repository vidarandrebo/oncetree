package oncetree

import "github.com/vidarandrebo/oncetree/protos"

type QSpec struct {
	numNodes int
}

func (q *QSpec) ReadAllQF(in *protos.ReadRequest, replies map[uint32]*protos.ReadAllResponse) (*protos.ReadAllResponse, bool) {
	return nil, false
}
