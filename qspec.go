package oncetree

import "github.com/vidarandrebo/oncetree/protos/keyvaluestorageprotos"

type FDQSpec struct {
	NumNodes int
}

type QSpec struct {
	NumNodes int
}

func (q *QSpec) ReadAllQF(in *keyvaluestorageprotos.ReadRequest, replies map[uint32]*keyvaluestorageprotos.ReadAllResponse) (*keyvaluestorageprotos.ReadAllResponse, bool) {
	return nil, false
}
