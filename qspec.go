package oncetree

import "github.com/vidarandrebo/oncetree/protos/keyvaluestorageprotos"

type QSpec struct {
	numNodes int
}

func (q *QSpec) ReadAllQF(in *keyvaluestorageprotos.ReadRequest, replies map[uint32]*keyvaluestorageprotos.ReadAllResponse) (*keyvaluestorageprotos.ReadAllResponse, bool) {
	return nil, false
}
