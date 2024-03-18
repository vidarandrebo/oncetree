package sqspec

import kvsprotos "github.com/vidarandrebo/oncetree/protos/keyvaluestorage"

type QSpec struct {
	NumNodes int
}

func (q *QSpec) ReadAllQF(in *kvsprotos.ReadRequest, replies map[uint32]*kvsprotos.ReadAllResponse) (*kvsprotos.ReadAllResponse, bool) {
	if len(replies) < q.NumNodes {
		return nil, false
	}
	values := make(map[string]int64)
	// merges the response maps into one maps. The maps in each reply only contains one key-value pair, but this is needed to have the same response type for the individaul response and the quorum call.
	for _, reply := range replies {
		for id, value := range reply.Value {
			values[id] = value
		}
	}
	return &kvsprotos.ReadAllResponse{Value: values}, true
}
