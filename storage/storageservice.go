package storage

import (
	"context"
	"log"
	"reflect"

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
	storage        KeyValueStorage
	logger         *log.Logger
	timestamp      *mutex.RWMutex[int64]
	nodeManager    *nodemanager.NodeManager
	eventBus       *eventbus.EventBus
	configProvider gorumsprovider.StorageConfigProvider
}

func NewStorageService(id string, logger *log.Logger, nodeManager *nodemanager.NodeManager, eventBus *eventbus.EventBus, configProvider gorumsprovider.StorageConfigProvider) *StorageService {
	ss := &StorageService{
		id:             id,
		logger:         logger,
		storage:        *NewKeyValueStorage(),
		nodeManager:    nodeManager,
		timestamp:      mutex.New[int64](0),
		eventBus:       eventBus,
		configProvider: configProvider,
	}
	//	eventBus.RegisterHandler(reflect.TypeOf(nodemanager.NeighbourAddedEvent{}), func(e any) {
	//	})
	eventBus.RegisterHandler(reflect.TypeOf(nodemanager.NeighbourReadyEvent{}), func(e any) {
		if event, ok := e.(nodemanager.NeighbourReadyEvent); ok {
			ss.shareAll(event.NodeID)
		}
	})
	return ss
}

func (ss *StorageService) shareAll(nodeID string) {
	ss.logger.Printf("[StorageService] - sharing all values with %s", nodeID)
	neighbour, ok := ss.nodeManager.Neighbour(nodeID)
	if !ok {
		ss.logger.Printf("[StorageService] - did not find neighbour %v", nodeID)
		return
	}
	gorumsConfig := ss.configProvider.StorageConfig()
	node, ok := gorumsConfig.Node(neighbour.GorumsID)
	if !ok {
		ss.logger.Printf("[StorageService] - did not find node %d in gorums config", neighbour.GorumsID)
		return
	}
	for _, key := range ss.storage.Keys().Values() {
		tsRef := ss.timestamp.RLock()
		ts := *tsRef
		value, err := ss.storage.ReadValueExceptNode(nodeID, key)
		ss.timestamp.RUnlock(&tsRef)
		if err != nil {
			ss.logger.Println(err)
			continue
		}
		ctx, cancel := context.WithTimeout(context.Background(), consts.RPCContextTimeout)
		request := &kvsprotos.GossipMessage{
			NodeID:    ss.id,
			Key:       key,
			Value:     value,
			Timestamp: ts,
		}
		_, err = node.Gossip(ctx, request)
		if err != nil {
			ss.logger.Println(err)
		}
		cancel()
	}
	ss.logger.Printf("[StorageService] - completed sharing values with node %v", nodeID)
}

func (ss *StorageService) sendGossip(originID string, key int64, values map[string]int64, ts int64, writeID int64) {
	sent := false
	// ss.logger.Printf("[StorageService] - sendGossip, key = %d, origin = %s, writeID = %d", key, originID, writeID)
	gorumsConfig := ss.configProvider.StorageConfig()
	for _, gorumsNode := range gorumsConfig.Nodes() {
		nodeID, ok := ss.nodeManager.NodeID(gorumsNode.ID())
		if !ok {
			ss.logger.Println("[StorageService] - node lookup failed")
			continue
		}
		// skip returning to originID and sending to self.
		if nodeID == originID {
			// ss.logger.Printf("[StorageService] - node is origin, writeID = %d", writeID)
			continue
		}
		if nodeID == ss.id {
			ss.logger.Panicf("[StorageService] - node is self, writeID = %d", writeID)
		}
		// value, err := ss.storage.ReadValueExceptNode(nodeID, key)
		ctx, cancel := context.WithTimeout(context.Background(), consts.RPCContextTimeout)
		_, err := gorumsNode.Gossip(ctx, &kvsprotos.GossipMessage{
			NodeID:    ss.id,
			Key:       key,
			Value:     values[nodeID],
			Timestamp: ts,
			WriteID:   writeID,
		},
		)
		sent = true
		if err != nil {
			ss.logger.Panicf("[StorageService] - Gossip rpc of writeID = %d to nodeID = %s err = %v", writeID, nodeID, err)
		}
		cancel()
	}
	if sent == false {
		// ss.logger.Printf("[StorageService] - gossip rpc of writeID = %d, node is leaf-node", writeID)
	}
}

func (ss *StorageService) Write(ctx gorums.ServerCtx, request *kvsprotos.WriteRequest) (response *emptypb.Empty, err error) {
	// ss.logger.Printf("[StorageService] - write rpc, key = %d, value = %d, writeID = %d", request.GetKey(), request.GetValue(), request.GetWriteID())
	ts := ss.timestamp.Lock()
	*ts++
	writeTs := *ts
	ok := ss.storage.WriteValue(ss.id, request.GetKey(), request.GetValue(), writeTs)
	valuesToGossip, err := ss.storage.GetGossipValues(
		request.GetKey(),
		ss.nodeManager.NeighbourIDs(),
	)
	// both the write and read must happen while ts mutex is locked to avoid inconsistencies
	ss.timestamp.Unlock(&ts)

	if err != nil {
		ss.logger.Panicln(err)
		return &emptypb.Empty{}, nil
	}
	if ok && ss.hasValueToGossip(ss.id, valuesToGossip) {
		// only start gossip if write was successful
		ss.eventBus.PushTask(func() {
			ss.sendGossip(ss.id, request.GetKey(), valuesToGossip, writeTs, request.GetWriteID())
		})
	} else {
		ss.logger.Printf("write to key %v failed because existing value has higher timestamp", writeTs)
	}
	return &emptypb.Empty{}, nil
}

func (ss *StorageService) Read(ctx gorums.ServerCtx, request *kvsprotos.ReadRequest) (response *kvsprotos.ReadResponse, err error) {
	value, err := ss.storage.ReadValue(request.Key)
	if err != nil {
		return &kvsprotos.ReadResponse{Value: 0}, err
	}
	return &kvsprotos.ReadResponse{Value: value}, nil
}

func (ss *StorageService) ReadAll(ctx gorums.ServerCtx, request *kvsprotos.ReadRequest) (response *kvsprotos.ReadAllResponse, err error) {
	value, err := ss.storage.ReadValue(request.Key)
	if err != nil {
		return &kvsprotos.ReadAllResponse{Value: nil}, err
	}
	return &kvsprotos.ReadAllResponse{Value: map[string]int64{ss.id: value}}, nil
}

func (ss *StorageService) PrintState(ctx gorums.ServerCtx, request *emptypb.Empty) (response *emptypb.Empty, err error) {
	ss.logger.Println(ss.storage.data)
	return &emptypb.Empty{}, nil
}

func (ss *StorageService) Gossip(ctx gorums.ServerCtx, request *kvsprotos.GossipMessage) (response *emptypb.Empty, err error) {
	// ctx.Release()
	// ss.logger.Printf("[StorageService] - gossip rpc, writeID = %d, nodeID = %s received", request.GetWriteID(), request.GetNodeID())
	ts := ss.timestamp.Lock()
	*ts++
	writeTs := *ts
	updated := ss.storage.WriteValue(request.GetNodeID(), request.GetKey(), request.GetValue(), request.GetTimestamp())
	valuesToGossip, err := ss.storage.GetGossipValues(
		request.GetKey(),
		ss.nodeManager.NeighbourIDs(),
	)
	// both the write and read must happen while ts mutex is locked to avoid inconsistencies
	ss.timestamp.Unlock(&ts)

	if err != nil {
		ss.logger.Panicf("[StorageService]177 - %v", err)
		return &emptypb.Empty{}, nil
	}

	if updated && ss.hasValueToGossip(request.GetNodeID(), valuesToGossip) {
		// go ss.sendGossip(request.NodeID, request.GetKey(), valuesToGossip, writeTs)
		ss.eventBus.PushTask(func() {
			go ss.sendGossip(request.GetNodeID(), request.GetKey(), valuesToGossip, writeTs, request.GetWriteID())
		})
	}
	// ss.logger.Printf("[StorageService] - gossip rpc, writeID = %d, nodeID = %s handled\n", request.GetWriteID(), request.GetNodeID())
	return &emptypb.Empty{}, nil
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
