package oncetree

import "github.com/vidarandrebo/oncetree/protos/keyvaluestorage"

type FDQSpec struct {
	NumNodes int
}

type QSpec struct {
	NumNodes int
}

func (q *QSpec) ReadAllQF(in *keyvaluestorage.ReadRequest, replies map[uint32]*keyvaluestorage.ReadAllResponse) (*keyvaluestorage.ReadAllResponse, bool) {
	if len(replies) < q.NumNodes {
		return nil, false
	}
	values := make(map[string]int64)
	// merges the response maps into one map. The map in each reply only contains one key-value pair, but this is needed to have the same response type for the individaul response and the quorum call.
	for _, reply := range replies {
		for id, value := range reply.Value {
			values[id] = value
		}
	}
	return &keyvaluestorage.ReadAllResponse{Value: values}, true
}
