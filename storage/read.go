package storage

import (
	"errors"
	"log/slog"

	"github.com/relab/gorums"
	kvsprotos "github.com/vidarandrebo/oncetree/protos/keyvaluestorage"
)

func (ss *StorageService) Read(ctx gorums.ServerCtx, request *kvsprotos.ReadRequest) (*kvsprotos.ReadResponse, error) {
	ss.requestMetrics.CountRead()
	value, err := ss.storage.ReadValue(request.Key)
	if err != nil {
		return &kvsprotos.ReadResponse{Value: 0}, err
	}
	ss.logger.Debug("Read rpc",
		slog.Int64("key", request.GetKey()),
		slog.Int64("value", value))
	return &kvsprotos.ReadResponse{Value: value}, nil
}

// ReadLocal rpc is used for checking that local values are propagated as intended
func (ss *StorageService) ReadLocal(ctx gorums.ServerCtx, request *kvsprotos.ReadLocalRequest) (*kvsprotos.ReadResponse, error) {
	value, ok := ss.storage.ReadLocalValue(request.GetKey(), request.GetNodeID())
	if !ok {
		ss.logger.Debug("ReadLocal rpc", slog.Int64("value", -1))
		return &kvsprotos.ReadResponse{Value: 0}, errors.New("value not found")
	}
	ss.logger.Debug("ReadLocal rpc",
		slog.Int64("key", request.GetKey()),
		slog.Int64("value", value.Value))
	return &kvsprotos.ReadResponse{Value: value.Value}, nil
}

func (ss *StorageService) ReadAll(ctx gorums.ServerCtx, request *kvsprotos.ReadRequest) (*kvsprotos.ReadResponseWithID, error) {
	ss.requestMetrics.CountRead()
	value, err := ss.storage.ReadValue(request.Key)
	if err != nil {
		return &kvsprotos.ReadResponseWithID{
			Value: 0,
			ID:    ss.id,
		}, err
	}
	return &kvsprotos.ReadResponseWithID{
			Value: value,
			ID:    ss.id,
		},
		nil
}
