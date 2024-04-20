package storage

import (
	"context"
	"errors"
	"log/slog"
	"reflect"

	"github.com/vidarandrebo/oncetree/nodemanager/nmevents"

	"github.com/vidarandrebo/oncetree/gorumsprovider"

	"github.com/relab/gorums"
	"github.com/vidarandrebo/oncetree/concurrent/mutex"
	"github.com/vidarandrebo/oncetree/consts"
	"github.com/vidarandrebo/oncetree/eventbus"
	"github.com/vidarandrebo/oncetree/nodemanager"
	kvsprotos "github.com/vidarandrebo/oncetree/protos/keyvaluestorage"
	"google.golang.org/protobuf/types/known/emptypb"
)

type StorageService struct {
	id             string
	storage        *KeyValueStorage
	logger         *slog.Logger
	timestamp      *mutex.RWMutex[int64]
	nodeManager    *nodemanager.NodeManager
	eventBus       *eventbus.EventBus
	configProvider gorumsprovider.StorageConfigProvider
}

func NewStorageService(id string, logger *slog.Logger, nodeManager *nodemanager.NodeManager, eventBus *eventbus.EventBus, configProvider gorumsprovider.StorageConfigProvider) *StorageService {
	ss := &StorageService{
		id:             id,
		logger:         logger.With(slog.Group("node", slog.String("module", "storageservice"))),
		storage:        NewKeyValueStorage(),
		nodeManager:    nodeManager,
		timestamp:      mutex.New[int64](0),
		eventBus:       eventBus,
		configProvider: configProvider,
	}
	eventBus.RegisterHandler(reflect.TypeOf(nmevents.NeighbourReadyEvent{}), func(e any) {
		if event, ok := e.(nmevents.NeighbourReadyEvent); ok {
			ss.shareAll(event.NodeID)
		}
	})
	eventBus.RegisterHandler(reflect.TypeOf(nmevents.TreeRecoveredEvent{}), func(e any) {
		if event, ok := e.(nmevents.TreeRecoveredEvent); ok {
			logger.Info(event.FailedNodeID)
			ss.MergeKeys(event)
		}
	})
	return ss
}

func (ss *StorageService) MergeKeys(event nmevents.TreeRecoveredEvent) {
	//keySet := ss.storage.Keys()
	_, ok := ss.configProvider.StorageConfig()
	if !ok {
		ss.logger.Error("failed to get storage-config")
	}
	//	for _, key := range keySet.Values() {
	//		storedValue, ok := ss.storage.ReadLocalValue(key, event.FailedNodeID)
	//		ctx, cancel := context.WithTimeout(context.Background(), consts.RPCContextTimeout)
	//		response, err := cfg.Prepare(ctx, &kvsprotos.PrepareMessage{
	//			Key:          key,
	//			Ts:           storedValue.Timestamp,
	//			FailedNodeID: event.FailedNodeID,
	//		})
	//		cancel()
	//	}
}

func (ss *StorageService) shareAll(nodeID string) {
	ss.logger.Info("sharing all values", "id", nodeID)
	neighbour, ok := ss.nodeManager.Neighbour(nodeID)
	if !ok {
		ss.logger.Error("did not find neighbour", "id", nodeID)
		return
	}
	gorumsConfig, configExists := ss.configProvider.StorageConfig()
	if !configExists {
		ss.logger.Error("storageconfig does not exist",
			slog.String("fn", "ss.shareAll"))
		return
	}
	node, ok := gorumsConfig.Node(neighbour.GorumsID)
	if !ok {
		ss.logger.Error(
			"did not find node in gorums config",
			slog.Uint64("id", uint64(neighbour.GorumsID)))
		return
	}
	for _, key := range ss.storage.Keys().Values() {
		tsRef := ss.timestamp.RLock()
		ts := *tsRef
		value, err := ss.storage.ReadValueExceptNode(nodeID, key)
		// TODO - get local value
		ss.timestamp.RUnlock(&tsRef)
		if err != nil {
			ss.logger.Error(
				"failed to read from storage",
				slog.Int64("key", key),
				slog.Any("err", err))
			continue
		}
		ctx, cancel := context.WithTimeout(context.Background(), consts.RPCContextTimeout)
		request := &kvsprotos.GossipMessage{
			NodeID:       ss.id,
			Key:          key,
			AggValue:     value,
			AggTimestamp: ts,
			// insert local value
		}
		_, err = node.Gossip(ctx, request)
		if err != nil {
			ss.logger.Error(
				"gossip failed",
				slog.Any("err", err))
		}
		cancel()
	}
	ss.logger.Info(
		"completed sharing values",
		slog.String("id", nodeID))
}

// hasValueToGossip determines if any of the values stored will have to be gossiped
//
// This fn is needed to make the Gossip fn converge when the task-queue is bounded
func (ss *StorageService) hasValueToGossip(origin string, valuesToGossip map[string]int64) bool {
	for key := range valuesToGossip {
		if (key != ss.id) && (key != origin) {
			return true
		}
	}
	return false
}

func (ss *StorageService) sendGossip(originID string, key int64, values map[string]int64, ts int64, localValue TimestampedValue, writeID int64) {
	ss.logger.Debug(
		"sendGossip",
		slog.String("originID", originID),
		slog.Int64("key", key),
		slog.Int64("ts", ts),
		slog.Int64("writeID", writeID),
	)
	sent := false
	gorumsConfig, configExists := ss.configProvider.StorageConfig()
	if !configExists {
		ss.logger.Error("storageconfig does not exist",
			slog.String("fn", "ss.shareAll"))
		return
	}
	for _, gorumsNode := range gorumsConfig.Nodes() {
		nodeID, ok := ss.nodeManager.NodeID(gorumsNode.ID())
		if !ok {
			ss.logger.Error(
				"node lookup failed",
				slog.Uint64("gorumsID", uint64(gorumsNode.ID())),
			)
			continue
		}
		// skip returning to originID and sending to self.
		if nodeID == originID {
			ss.logger.Debug(
				"target is origin, skipping gossip",
				slog.String("id", nodeID))
			continue
		}
		if nodeID == ss.id {
			ss.logger.Error("target is self, gorums configuration error")
		}
		ctx, cancel := context.WithTimeout(context.Background(), consts.RPCContextTimeout)
		_, err := gorumsNode.Gossip(ctx, &kvsprotos.GossipMessage{
			NodeID:         ss.id,
			Key:            key,
			AggValue:       values[nodeID],
			AggTimestamp:   ts,
			LocalValue:     localValue.Value,
			LocalTimestamp: localValue.Timestamp,
			WriteID:        writeID,
		},
		)
		sent = true
		if err != nil {
			ss.logger.Error("sending of gossip message failed",
				slog.Any("err", err),
				slog.Int64("key", key),
				slog.String("nodeID", nodeID),
				slog.Int64("value", values[nodeID]))
		}
		cancel()
	}
	if sent == false {
		ss.logger.Debug("node is leaf node, no message send",
			slog.Any("key", key))
	}
}

// Write RPC is used to set the local value at a given key for a single node
func (ss *StorageService) Write(ctx gorums.ServerCtx, request *kvsprotos.WriteRequest) (*emptypb.Empty, error) {
	ss.logger.Debug("RPC Write", "key", request.GetKey(), "value", request.GetValue(), "writeID", request.GetWriteID())
	ts := ss.timestamp.Lock()
	*ts++
	writeTs := *ts
	ok := ss.storage.WriteValue(ss.id, request.GetKey(), request.GetValue(), writeTs)

	// does not user storage.ReadLocalValue, since the aggregated value IS the local value in this case, and self's local value cannot be found using ReadLocalValue
	localValue, readValueErr := ss.storage.ReadValueFromNode(ss.id, request.GetKey())
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
		ss.logger.Debug("node does not have local value for this key",
			slog.Int64("key", request.GetKey()))
	}

	if ok && ss.hasValueToGossip(ss.id, valuesToGossip) {
		// only start gossip if write was successful
		ss.sendGossip(ss.id, request.GetKey(), valuesToGossip, writeTs, TimestampedValue{Value: localValue, Timestamp: writeTs}, request.GetWriteID())
	} else {
		ss.logger.Warn(
			"write failed because existing value has higher timestamp",
			slog.Int64("key", request.GetKey()),
			slog.Int64("ts", writeTs))
	}
	return &emptypb.Empty{}, nil
}

func (ss *StorageService) Read(ctx gorums.ServerCtx, request *kvsprotos.ReadRequest) (*kvsprotos.ReadResponse, error) {
	value, err := ss.storage.ReadValue(request.Key)
	if err != nil {
		return &kvsprotos.ReadResponse{Value: 0}, err
	}
	return &kvsprotos.ReadResponse{Value: value}, nil
}

// ReadLocal rpc is used for checking that local values are propagated as intended
func (ss *StorageService) ReadLocal(ctx gorums.ServerCtx, request *kvsprotos.ReadLocalRequest) (*kvsprotos.ReadResponse, error) {
	ss.logger.Debug("ReadLocal rpc",
		slog.Int64("key", request.GetKey()),
		slog.String("nodeID", request.GetNodeID()))
	value, ok := ss.storage.ReadLocalValue(request.GetKey(), request.GetNodeID())
	if !ok {
		return &kvsprotos.ReadResponse{Value: 0}, errors.New("value not found")
	}
	return &kvsprotos.ReadResponse{Value: value.Value}, nil
}

func (ss *StorageService) ReadAll(ctx gorums.ServerCtx, request *kvsprotos.ReadRequest) (response *kvsprotos.ReadAllResponse, err error) {
	value, err := ss.storage.ReadValue(request.Key)
	if err != nil {
		return &kvsprotos.ReadAllResponse{Value: nil}, err
	}
	return &kvsprotos.ReadAllResponse{Value: map[string]int64{ss.id: value}}, nil
}

func (ss *StorageService) PrintState(ctx gorums.ServerCtx, request *emptypb.Empty) (*emptypb.Empty, error) {
	// ss.logger.Println(ss.storage.data)
	return &emptypb.Empty{}, nil
}

func (ss *StorageService) Gossip(ctx gorums.ServerCtx, request *kvsprotos.GossipMessage) (*emptypb.Empty, error) {
	ss.logger.Debug(
		"RPC Gossip",
		slog.Int64("key", request.GetKey()),
		slog.Int64("value", request.GetAggValue()),
		slog.Int64("ts", request.GetAggTimestamp()),
		slog.String("nodeID", request.GetNodeID()),
		slog.Int64("writeID", request.GetWriteID()),
	)
	ctx.Release()
	ts := ss.timestamp.Lock()
	*ts++
	writeTs := *ts
	updated := ss.storage.WriteValue(request.GetNodeID(), request.GetKey(), request.GetAggValue(), request.GetAggTimestamp())
	wroteLocal := ss.storage.WriteLocalValue(request.GetNodeID(), request.GetKey(), request.GetLocalValue(), request.GetLocalTimestamp())
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

	if updated && ss.hasValueToGossip(request.GetNodeID(), valuesToGossip) {
		go ss.sendGossip(request.NodeID, request.GetKey(), valuesToGossip, writeTs, localValue, request.GetWriteID())
	}
	return &emptypb.Empty{}, nil
}

func (ss *StorageService) Prepare(ctx gorums.ServerCtx, request *kvsprotos.PrepareMessage) (*kvsprotos.PromiseMessage, error) {
	ss.logger.Info("Prepare RPC",
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

func (ss *StorageService) Accept(ctx gorums.ServerCtx, request *kvsprotos.AcceptMessage) (*kvsprotos.LearnMessage, error) {
	//TODO implement me
	panic("implement me")
}
