package storage

import (
	"log/slog"

	"github.com/relab/gorums"
	kvsprotos "github.com/vidarandrebo/oncetree/protos/keyvaluestorage"
)

func (ss *StorageService) Accept(ctx gorums.ServerCtx, request *kvsprotos.AcceptMessage) (*kvsprotos.LearnMessage, error) {
	ss.logger.Info("ss.Accept RPC",
		slog.String("failedNodeID", request.GetFailedNodeID()),
		slog.String("id", request.GetNodeID()),
		slog.Int64("key", request.GetKey()),
		slog.Int64("localValue", request.GetLocalValue()),
		slog.Int64("aggValue", request.GetAggValue()))

	ss.storage.ChangeKeyValueOwnership(
		request.GetKey(),
		request.GetAggValue(),
		request.GetLocalValue(),
		request.GetTimestamp(),
		request.GetFailedNodeID(),
		request.GetNodeID(),
	)
	// TODO - send gossip
	return &kvsprotos.LearnMessage{OK: true}, nil
}
