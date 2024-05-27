package storage

import (
	"log/slog"

	"github.com/relab/gorums"
	kvsprotos "github.com/vidarandrebo/oncetree/protos/keyvaluestorage"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (ss *StorageService) Gossip(ctx gorums.ServerCtx, request *kvsprotos.GossipMessage) (*emptypb.Empty, error) {
	ss.requestMetrics.CountGossip()
	ss.logger.Debug(
		"RPC Gossip",
		slog.Int64("key", request.GetKey()),
		slog.Int64("value", request.GetAggValue()),
		slog.Int64("ts", request.GetAggTimestamp()),
		slog.String("nodeID", request.GetNodeID()),
	)
	// ctx.Release()
	ts := ss.timestamp.Lock()
	*ts++
	writeTs := *ts
	updated := ss.storage.WriteValue(request.GetKey(), request.GetAggValue(), request.GetAggTimestamp(), request.GetNodeID())
	wroteLocal := ss.storage.WriteLocalValue(request.GetKey(), request.GetLocalValue(), request.GetLocalTimestamp(), request.GetNodeID())
	if wroteLocal {
		ss.logger.Debug("wrote local value",
			slog.Int64("key", request.GetKey()),
			slog.Int64("value", request.GetLocalValue()))
	} else {
		ss.logger.Debug("did not receive local value",
			slog.Int64("key", request.GetKey()))
	}
	localValue, hasLocal := ss.storage.ReadLocalValue(request.GetKey(), ss.id)
	valuesToGossip, gossipValueErr := ss.storage.GossipValues(
		request.GetKey(),
		ss.nodeManager.NeighbourIDs(),
	)
	// both the write and read must happen while ts mutex is locked to avoid inconsistencies
	ss.timestamp.Unlock(&ts)

	if gossipValueErr != nil {
		ss.logger.Warn("failed to retrieve values to gossip",
			slog.Any("err", gossipValueErr),
			slog.Int64("key", request.GetKey()),
			slog.String("id", request.GetNodeID()))
		return &emptypb.Empty{}, nil
	}
	if !hasLocal {
		ss.logger.Debug("node does not have local value for this key",
			slog.Int64("key", request.GetKey()),
			slog.String("id", request.GetNodeID()))
	}

	if updated {
		ss.sendGossip(request.NodeID, request.GetKey(), valuesToGossip, writeTs, localValue)
		ss.logger.Debug("value updated, stored value is older",
			slog.String("nodeID", request.GetNodeID()),
			slog.Int64("key", request.GetKey()))
	} else {
		ss.logger.Debug("value not updated, stored value is newer",
			slog.String("nodeID", request.GetNodeID()),
			slog.Int64("key", request.GetKey()))
	}
	return &emptypb.Empty{}, nil
}
