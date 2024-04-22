package storage

import (
	"log/slog"

	"github.com/relab/gorums"
	kvsprotos "github.com/vidarandrebo/oncetree/protos/keyvaluestorage"
)

func (ss *StorageService) Prepare(ctx gorums.ServerCtx, request *kvsprotos.PrepareMessage) (*kvsprotos.PromiseMessage, error) {
	ss.logger.Info("Prepare RPC",
		slog.Int64("key", request.GetKey()),
		slog.String("failedNodeID", request.GetFailedNodeID()))
	value, ok := ss.storage.ReadLocalValue(request.GetKey(), request.GetFailedNodeID())
	if !ok {
		// no local value stored
		return &kvsprotos.PromiseMessage{
			OK:    true,
			Value: 0,
			Ts:    0,
		}, nil
	}
	if value.Timestamp <= request.GetTs() {
		// older or same value stored
		return &kvsprotos.PromiseMessage{
			OK:    true,
			Value: 0,
			Ts:    0,
		}, nil
	}
	// newer value stored
	return &kvsprotos.PromiseMessage{
		OK:    false,
		Value: value.Value,
		Ts:    value.Timestamp,
	}, nil
}
