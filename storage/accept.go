package storage

import (
	"log/slog"

	"github.com/relab/gorums"
	kvsprotos "github.com/vidarandrebo/oncetree/protos/keyvaluestorage"
)

func (ss *StorageService) Accept(ctx gorums.ServerCtx, request *kvsprotos.AcceptMessage) (*kvsprotos.LearnMessage, error) {
	ss.requestMetrics.CountAccept()
	ss.logger.Debug("ss.Accept RPC",
		slog.String("failedNodeID", request.GetFailedNodeID()),
		slog.String("id", request.GetNodeID()),
		slog.Int64("key", request.GetKey()),
		slog.Int64("localValue", request.GetLocalValue()),
		slog.Int64("aggValue", request.GetAggValue()))

	tsRef := ss.timestamp.Lock()
	*tsRef++
	ts := *tsRef

	ss.storage.ChangeKeyValueOwnership(
		request.GetKey(),
		request.GetAggValue(),
		request.GetLocalValue(),
		request.GetTimestamp(),
		request.GetFailedNodeID(),
		request.GetNodeID(),
	)
	ss.storage.RemoveExclusionsFromKey(request.GetKey())

	localValue, err := ss.storage.ReadValueFromNode(request.GetKey(), ss.id)
	if err != nil {
		ss.logger.Warn("failed to get local value to gossip",
			slog.Any("err", err),
			slog.Int64("key", request.GetKey()))
		localValue = TimestampedValue{Value: 0, Timestamp: 0}
	}

	valuesToGossip, err := ss.storage.GossipValues(request.GetKey(), ss.nodeManager.NeighbourIDs())
	if err != nil {
		ss.logger.Error("failed to retrieve values to gossip",
			slog.Int64("key", request.GetKey()))
	}
	ss.timestamp.Unlock(&tsRef)

	go ss.sendGossip(request.GetNodeID(), request.GetKey(), valuesToGossip, ts, localValue)
	return &kvsprotos.LearnMessage{
		OK: true,
	}, nil
}
