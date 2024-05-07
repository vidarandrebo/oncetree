package storage

import (
	"context"
	"log/slog"
	"reflect"

	"github.com/vidarandrebo/oncetree/storage/sevents"

	"github.com/vidarandrebo/oncetree/nodemanager/nmevents"

	"github.com/vidarandrebo/oncetree/gorumsprovider"

	"github.com/relab/gorums"
	"github.com/vidarandrebo/oncetree/concurrent/maps"
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
			ss.SendPrepare(event)
		}
	})
	return ss
}

func (ss *StorageService) SendPrepare(event nmevents.TreeRecoveredEvent) {
	keySet := ss.storage.Keys()
	cfg, ok, _ := ss.configProvider.StorageConfig()
	if !ok {
		ss.logger.Error("failed to get storage-config")
	}
	for _, key := range keySet.Values() {
		storedValue, ok := ss.storage.ReadLocalValue(key, event.FailedNodeID)
		if !ok {
			ss.logger.Info("storage does not contain value",
				slog.Int64("key", key))
			continue
			// TODO, might need to actually send message here
		}
		ctx, cancel := context.WithTimeout(context.Background(), consts.RPCContextTimeout)
		response, err := cfg.Prepare(ctx, &kvsprotos.PrepareMessage{
			NodeID:       ss.id,
			Key:          key,
			Ts:           storedValue.Timestamp,
			FailedNodeID: event.FailedNodeID,
		})
		cancel()
		if err != nil {
			panic("failed to get response")
		}
		if response.GetOK() {
			ss.logger.Info("node has latest ts",
				slog.Int64("key", key),
				slog.Int64("ts", storedValue.Timestamp))
		} else {
			// this case is extremely unlikely as the failing node will have had to transmit an update to only a
			// subset of its neighbours while failing
			ss.logger.Warn("remote has latest ts",
				slog.Int64("key", key),
				slog.Int64("ts", response.GetTs()))
			ss.storage.WriteLocalValue(key, response.GetValue(), response.GetTs(), event.FailedNodeID)
		}
		// TODO - there is a case where a failed node has partially propagated a new key.
		// there should be a way to ensure that the values of all keys are transferred into leaders ownership,
		// not just the keys the leader knows of.
	}
	ss.SendAccept(event.FailedNodeID)
}

func (ss *StorageService) SendAccept(failedNodeID string) {
	ss.logger.Info("sending Accept messages",
		slog.String("fn", "ss.SendAccept"))
	keySet := ss.storage.Keys()
	cfg, ok, _ := ss.configProvider.StorageConfig()
	if !ok {
		ss.logger.Error("failed to get storage-config")
	}
	for _, key := range keySet.Values() {
		ts := ss.timestamp.Lock()
		*ts++
		writeTs := *ts
		localValue, readValueErr := ss.storage.ReadValueFromNode(key, ss.id)
		if readValueErr != nil {
			ss.logger.Info("node does not contain value",
				slog.Int64("key", key),
				slog.Any("err", readValueErr))
			localValue = TimestampedValue{Value: 0, Timestamp: 0}
		}
		oldLocal, ok := ss.storage.ReadLocalValue(key, failedNodeID)
		if !ok && (readValueErr != nil) {
			ss.logger.Info("failed node and self does not store value for key",
				slog.String("id", failedNodeID),
				slog.Int64("key", key))
			ss.timestamp.Unlock(&ts)
			continue
		} else if !ok {
			oldLocal = TimestampedValue{Value: 0, Timestamp: 0}
		}
		newLocalValue := oldLocal.Value + localValue.Value

		ok = ss.storage.WriteValue(key, newLocalValue, writeTs, ss.id)
		if !ok {
			ss.logger.Error("failed to write local new value",
				slog.Int64("value", newLocalValue),
				slog.Int64("key", key))
		}
		ss.storage.DeleteAgg(key, failedNodeID)
		ss.storage.DeleteLocal(key, failedNodeID)

		// does not user storage.ReadLocalValue, since the aggregated value IS the local value in this case, and self's local value cannot be found using ReadLocalValue
		// both the write and read must happen while ts mutex is locked to avoid inconsistencies
		valuesToGossip, gossipValueErr := ss.storage.GossipValues(
			key,
			ss.nodeManager.NeighbourIDs(),
		)
		ss.timestamp.Unlock(&ts)

		if gossipValueErr != nil {
			ss.logger.Error("failed to get valus to gossip")
		}

		valuesToGossipSafe := maps.FromMap(valuesToGossip)
		ctx, cancel := context.WithTimeout(context.Background(), consts.RPCContextTimeout)
		response, err := cfg.Accept(ctx, &kvsprotos.AcceptMessage{
			NodeID:       ss.id,
			Key:          key,
			AggValue:     -1,
			LocalValue:   newLocalValue,
			FailedNodeID: failedNodeID,
			Timestamp:    writeTs,
		}, func(request *kvsprotos.AcceptMessage, gorumsID uint32) *kvsprotos.AcceptMessage {
			alteredRequest := kvsprotos.AcceptMessage{
				Key:          request.GetKey(),
				LocalValue:   request.GetLocalValue(),
				Timestamp:    request.GetTimestamp(),
				NodeID:       request.GetNodeID(),
				FailedNodeID: request.GetFailedNodeID(),
			}
			id, ok := ss.nodeManager.NodeID(gorumsID)
			if !ok {
				return request
			}
			aggValue, ok := valuesToGossipSafe.Get(id)
			if !ok {
				return request
			}
			alteredRequest.AggValue = aggValue
			return &alteredRequest
		})
		cancel()
		if err != nil {
			ss.logger.Error("sending accept failed",
				slog.Any("err", err))
		}
		ss.logger.Info("got response from accept",
			slog.Bool("ok", response.GetOK()))
	}
}

func (ss *StorageService) shareAll(nodeID string) {
	ss.logger.Info("sharing all values", "id", nodeID)
	neighbour, ok := ss.nodeManager.Neighbour(nodeID)
	if !ok {
		ss.logger.Error("did not find neighbour", "id", nodeID)
		return
	}
	gorumsConfig, configExists, _ := ss.configProvider.StorageConfig()
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
		localValue, err := ss.storage.ReadValueFromNode(key, ss.id)
		if err != nil {
			// not an error state in this case, just means that no local value exists for key
			localValue = TimestampedValue{Value: 0, Timestamp: 0}
		}
		aggValue, err := ss.storage.ReadValueExceptNode(key, nodeID)
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
			NodeID:         ss.id,
			Key:            key,
			AggValue:       aggValue,
			AggTimestamp:   ts,
			LocalValue:     localValue.Value,
			LocalTimestamp: localValue.Timestamp,
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

func (ss *StorageService) sendGossip(originID string, key int64, values map[string]int64, ts int64, localValue TimestampedValue) {
	ss.logger.Debug(
		"sendGossip",
		slog.String("originID", originID),
		slog.Int64("key", key),
		slog.Int64("ts", ts),
	)
	sent := false
	gorumsConfig, configExists, _ := ss.configProvider.StorageConfig()
	if !configExists {
		ss.logger.Error("storageconfig does not exist",
			slog.String("fn", "ss.shareAll"))
		return
	}
	for _, gorumsNode := range gorumsConfig.Nodes() {
		nodeID, ok := ss.nodeManager.NodeID(gorumsNode.ID())
		if ss.nodeManager.IsBlacklisted(nodeID) {
			ss.logger.Warn("skipping node, blacklisted",
				slog.String("nodeID", nodeID))
			continue
		}
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

		// ctx, cancel := context.WithTimeout(context.Background(), consts.RPCContextTimeout)
		ctx := context.Background()
		_, err := gorumsNode.Gossip(ctx, &kvsprotos.GossipMessage{
			NodeID:         ss.id,
			Key:            key,
			AggValue:       values[nodeID],
			AggTimestamp:   ts,
			LocalValue:     localValue.Value,
			LocalTimestamp: localValue.Timestamp,
		},
		)
		// cancel()
		sent = true
		if err != nil {
			ss.logger.Error("sending of gossip message failed",
				slog.Any("err", err),
				slog.Int64("key", key),
				slog.String("nodeID", nodeID),
				slog.Int64("value", values[nodeID]))
			ss.eventBus.PushEvent(sevents.NewGossipFailedEvent(nodeID))
		}
	}
	if sent == false {
		ss.logger.Debug("node is leaf node, no message send",
			slog.Any("key", key))
	}
}

func (ss *StorageService) PrintState(ctx gorums.ServerCtx, request *emptypb.Empty) (*emptypb.Empty, error) {
	// ss.logger.Println(ss.storage.data)
	return &emptypb.Empty{}, nil
}
