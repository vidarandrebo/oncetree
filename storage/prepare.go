package storage

import (
	"log/slog"

	"github.com/relab/gorums"
	kvsprotos "github.com/vidarandrebo/oncetree/protos/keyvaluestorage"
)

func (ss *StorageService) Prepare(ctx gorums.ServerCtx, request *kvsprotos.PrepareMessage) (*kvsprotos.PromiseMessage, error) {
	ss.requestMetrics.CountPrepare()
	ss.logger.Debug("Prepare RPC",
		slog.Int64("key", request.GetKey()),
		slog.String("failedNodeID", request.GetFailedNodeID()))
	failedNodeLocal, hasFailedLocalValue := ss.storage.ReadLocalValue(request.GetKey(), request.GetFailedNodeID())
	tsRef := ss.timestamp.RLock()
	aggValue, err := ss.storage.GossipValue(request.GetKey(), request.GetNodeID())
	ts := *tsRef
	ss.timestamp.RUnlock(&tsRef)
	if err != nil {
		aggValue = 0
	}
	local, err := ss.storage.ReadValueFromNode(request.GetKey(), ss.id)
	if err != nil {
		local = TimestampedValue{
			Value:     0,
			Timestamp: 0,
		}
	}
	promise := &kvsprotos.PromiseMessage{
		OK:                   true,
		AggValue:             aggValue,
		AggTimestamp:         ts,
		LocalValue:           local.Value,
		LocalTimestamp:       local.Timestamp,
		FailedLocalValue:     failedNodeLocal.Value,
		FailedLocalTimestamp: failedNodeLocal.Timestamp,
		NodeID:               ss.id,
	}
	if !hasFailedLocalValue {
		// no failedNodeLocal stored
		return promise, nil
	}
	if failedNodeLocal.Timestamp <= request.GetTs() {
		// older or same failedNodeLocal stored
		return promise, nil
	}

	// newer failedNodeLocal stored
	promise.OK = false
	return promise, nil
}
