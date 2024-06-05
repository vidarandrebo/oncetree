package sqspec

import kvsprotos "github.com/vidarandrebo/oncetree/protos/keyvaluestorage"

type QSpec struct {
	NumNodes int
}

func (q *QSpec) ReadAllQF(in *kvsprotos.ReadRequest, replies map[uint32]*kvsprotos.ReadResponseWithID) (*kvsprotos.ReadResponses, bool) {
	if len(replies) < q.NumNodes {
		return nil, false
	}
	values := make(map[string]int64)
	for _, reply := range replies {
		values[reply.ID] = reply.Value
	}
	return &kvsprotos.ReadResponses{Value: values}, true
}

func (q *QSpec) PrepareQF(in *kvsprotos.PrepareMessage, replies map[uint32]*kvsprotos.PromiseMessage) (*kvsprotos.PromiseMessages, bool) {
	if len(replies) < q.NumNodes {
		return nil, false
	}
	response := kvsprotos.PromiseMessages{
		OK:                   true,
		FailedLocalTimestamp: 0,
		FailedLocalValue:     0,
		Values:               make(map[string]*kvsprotos.PromiseMessage),
	}
	for _, reply := range replies {
		response.Values[reply.NodeID] = reply
		// OK means the data that was provided in prepare was up to date
		if reply.GetOK() {
			continue
		}
		// we overwrite the combined result if we have newer data in reply
		if reply.GetFailedLocalTimestamp() > response.GetFailedLocalTimestamp() {
			response.FailedLocalValue = reply.GetFailedLocalValue()
			response.FailedLocalTimestamp = reply.GetFailedLocalTimestamp()
			response.OK = false
		}

	}
	return &response, true
}

func (q *QSpec) AcceptQF(in *kvsprotos.AcceptMessage, replies map[uint32]*kvsprotos.LearnMessage) (*kvsprotos.LearnMessage, bool) {
	if len(replies) != q.NumNodes {
		return nil, false
	}

	learn := &kvsprotos.LearnMessage{
		OK: true,
	}
	for _, reply := range replies {
		if !reply.OK {
			learn.OK = false
			return learn, true
		}
	}
	return learn, true
}
