package storage

import (
	"context"
	"log"
	"reflect"
	"sync"

	"github.com/relab/gorums"
	"github.com/vidarandrebo/oncetree/concurrent/mutex"
	"github.com/vidarandrebo/oncetree/consts"
	"github.com/vidarandrebo/oncetree/eventbus"
	"github.com/vidarandrebo/oncetree/nodemanager"
	kvsprotos "github.com/vidarandrebo/oncetree/protos/keyvaluestorage"
	"google.golang.org/protobuf/types/known/emptypb"
)

type StorageService struct {
	id            string
	storage       KeyValueStorage
	logger        *log.Logger
	mut           sync.Mutex
	timestamp     *mutex.RWMutex[int64]
	gorumsConfig  *kvsprotos.Configuration
	gorumsManager *kvsprotos.Manager
	nodeManager   *nodemanager.NodeManager
	eventBus      *eventbus.EventBus
}

func NewStorageService(id string, logger *log.Logger, nodeManager *nodemanager.NodeManager, gorumsManager *kvsprotos.Manager, eventBus *eventbus.EventBus) *StorageService {
	ss := &StorageService{
		id:            id,
		logger:        logger,
		storage:       *NewKeyValueStorage(),
		gorumsManager: gorumsManager,
		nodeManager:   nodeManager,
		timestamp:     mutex.New[int64](0),
		eventBus:      eventBus,
	}
	eventBus.RegisterHandler(reflect.TypeOf(nodemanager.NeighbourAddedEvent{}), func(e any) {
		//	err := ss.SetNodesFromManager()
		//	if err != nil {
		//		ss.logger.Println(err)
		//	}
		//	if event, ok := e.(nodemanager.NeighbourAddedEvent); ok {
		//		ss.shareAll(event.NodeID)
		//	}
	})
	return ss
}

func (ss *StorageService) SetNodesFromManager() error {
	gorumsNeighbourMap := ss.nodeManager.GorumsNeighbourMap()
	cfg, err := ss.gorumsManager.NewConfiguration(
		&QSpec{
			NumNodes: len(gorumsNeighbourMap),
		},
		gorums.WithNodeMap(
			gorumsNeighbourMap,
		),
	)
	if err != nil {
		return err
	}
	ss.gorumsConfig = cfg
	return nil
}

func (ss *StorageService) shareAll(nodeID string) {
	tsRef := ss.timestamp.Lock()
	*tsRef++     // increment before sending
	ts := *tsRef // make sure all messages has same ts
	ss.timestamp.Unlock(&tsRef)
	neighbour, ok := ss.nodeManager.Neighbour(nodeID)
	if !ok {
		return
	}
	node, ok := ss.gorumsConfig.Node(neighbour.GorumsID)
	if !ok {
		return
	}
	for _, key := range ss.storage.Keys().Values() {
		ss.logger.Println(key)
		ss.logger.Println(node)
		ss.logger.Println(ts)
	}
}

func (ss *StorageService) sendGossip(originID string, key int64) {
	tsRef := ss.timestamp.Lock()
	*tsRef++     // increment before sending
	ts := *tsRef // make sure all messages has same ts
	ss.timestamp.Unlock(&tsRef)
	for _, node := range ss.gorumsConfig.Nodes() {
		nodeID, err := ss.nodeManager.ResolveNodeIDFromAddress(node.Address())
		if err != nil {
			continue
		}
		// skip returning to originID and sending to self.
		if nodeID == originID || nodeID == ss.id {
			continue
		}
		value, err := ss.storage.ReadValueExceptNode(nodeID, key)
		if err != nil {
			return
		}
		ctx, cancel := context.WithTimeout(context.Background(), consts.RPCContextTimeout)
		ctx.Done()
		_, err = node.Gossip(ctx, &kvsprotos.GossipMessage{NodeID: ss.id, Key: key, Value: value, Timestamp: ts})
		ss.logger.Println(err)
		cancel()
	}
}

func (ss *StorageService) Write(ctx gorums.ServerCtx, request *kvsprotos.WriteRequest) (response *emptypb.Empty, err error) {
	ss.logger.Println("got writerequest")
	ts := ss.timestamp.Lock()
	*ts++
	ss.storage.WriteValue(ss.id, request.GetKey(), request.GetValue(), *ts)
	ss.timestamp.Unlock(&ts)
	go func() {
		ss.sendGossip(ss.id, request.GetKey())
	}()
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
	ss.logger.Println(ss.storage)
	return &emptypb.Empty{}, nil
}

func (ss *StorageService) Gossip(ctx gorums.ServerCtx, request *kvsprotos.GossipMessage) (response *emptypb.Empty, err error) {
	ss.logger.Printf("received gossip %v", request)
	ss.mut.Lock()
	updated := ss.storage.WriteValue(request.GetNodeID(), request.GetKey(), request.GetValue(), request.GetTimestamp())
	ss.mut.Unlock()
	if updated {
		go func() {
			ss.sendGossip(request.NodeID, request.GetKey())
		}()
	}
	return &emptypb.Empty{}, nil
}
