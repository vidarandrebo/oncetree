package storage

import (
	"log/slog"

	"github.com/relab/gorums"
	kvsprotos "github.com/vidarandrebo/oncetree/protos/keyvaluestorage"
	"google.golang.org/protobuf/types/known/emptypb"
)

// Write RPC is used to set the local value at a given key for a single node
func (ss *StorageService) Write(ctx gorums.ServerCtx, request *kvsprotos.WriteRequest) (*emptypb.Empty, error) {
	ss.requestMetrics.CountWrite()
	ss.logger.Debug("RPC Write",
		slog.Int64("key", request.GetKey()),
		slog.Int64("value", request.GetValue()),
	)
	ts := ss.timestamp.Lock()
	*ts++
	writeTs := *ts
	ok := ss.storage.WriteValue(request.GetKey(), request.GetValue(), writeTs, ss.id)

	// does not user storage.ReadLocalValue, since the aggregated value IS the local value in this case, and self's local value cannot be found using ReadLocalValue
	localValue, readValueErr := ss.storage.ReadValueFromNode(request.GetKey(), ss.id)
	valuesToGossip, gossipValueErr := ss.storage.GossipValues(
		request.GetKey(),
		ss.nodeManager.NeighbourIDs(),
	)
	// both the write and read must happen while ts mutex is locked to avoid inconsistencies
	ss.timestamp.Unlock(&ts)

	if gossipValueErr != nil {
		ss.logger.Warn("failed to retrieve values to gossip",
			slog.Any("err", gossipValueErr),
			slog.Int64("key", request.GetKey()))
		return &emptypb.Empty{}, nil
	}
	if readValueErr != nil {
		localValue = TimestampedValue{Value: 0, Timestamp: 0}
		ss.logger.Debug("node does not have local value for this key",
			slog.Int64("key", request.GetKey()))
	}

	if ok {
		// only start gossip if write was successful
		ss.sendGossip(
			ss.id,
			request.GetKey(),
			valuesToGossip,
			writeTs,
			TimestampedValue{Value: localValue.Value, Timestamp: writeTs},
		)
	} else {
		ss.logger.Warn(
			"write failed because existing value has higher timestamp",
			slog.Int64("key", request.GetKey()),
			slog.Int64("ts", writeTs))
	}
	return &emptypb.Empty{}, nil
}
