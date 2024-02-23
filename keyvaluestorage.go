package oncetree

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/relab/gorums"
	kvsprotos "github.com/vidarandrebo/oncetree/protos/keyvaluestorageprotos"
	"google.golang.org/protobuf/types/known/emptypb"
)

type KeyValueStorageService struct {
	id            string
	storage       KeyValueStorage
	logger        *log.Logger
	mut           sync.Mutex
	timestamp     int64
	gorumsConfig  *kvsprotos.Configuration
	gorumsManager *kvsprotos.Manager
	nodeManager   *NodeManager
}

func NewKeyValueStorageService(id string, logger *log.Logger) *KeyValueStorageService {
	return &KeyValueStorageService{
		id:      id,
		logger:  logger,
		storage: *NewKeyValueStorage(),
	}
}

func (kvss *KeyValueStorageService) sendGossip(originID string, key int64) {
	kvss.mut.Lock()
	kvss.timestamp++
	ts := kvss.timestamp // make sure all messages has same ts
	kvss.mut.Unlock()
	for _, node := range kvss.gorumsConfig.Nodes() {
		nodeID, err := kvss.nodeManager.resolveNodeIDFromAddress(node.Address())
		if err != nil {
			continue
		}
		// skip returning to originID and sending to self.
		if nodeID == originID || nodeID == kvss.id {
			continue
		}
		value, err := kvss.storage.ReadValueExceptNode(nodeID, key)
		// TODO handle err
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		ctx.Done()
		_, err = node.Gossip(ctx, &kvsprotos.GossipMessage{NodeID: kvss.id, Key: key, Value: value, Timestamp: ts})
		// TODO handle err
		cancel()
	}
}

func (kvss *KeyValueStorageService) Write(ctx gorums.ServerCtx, request *kvsprotos.WriteRequest) (response *emptypb.Empty, err error) {
	kvss.mut.Lock()
	kvss.timestamp++
	kvss.storage.WriteValue(kvss.id, request.GetKey(), request.GetValue(), kvss.timestamp)
	kvss.mut.Unlock()
	go func() {
		kvss.sendGossip(kvss.id, request.GetKey())
	}()
	return &emptypb.Empty{}, nil
}

func (kvss *KeyValueStorageService) Read(ctx gorums.ServerCtx, request *kvsprotos.ReadRequest) (response *kvsprotos.ReadResponse, err error) {
	value, err := kvss.storage.ReadValue(request.Key)
	if err != nil {
		return &kvsprotos.ReadResponse{Value: 0}, err
	}
	return &kvsprotos.ReadResponse{Value: value}, nil
}

func (kvss *KeyValueStorageService) ReadAll(ctx gorums.ServerCtx, request *kvsprotos.ReadRequest) (response *kvsprotos.ReadAllResponse, err error) {
	value, err := kvss.storage.ReadValue(request.Key)
	if err != nil {
		return &kvsprotos.ReadAllResponse{Value: nil}, err
	}
	return &kvsprotos.ReadAllResponse{Value: map[string]int64{kvss.id: value}}, nil
}

func (kvss *KeyValueStorageService) PrintState(ctx gorums.ServerCtx, request *emptypb.Empty) (response *emptypb.Empty, err error) {
	kvss.logger.Println(kvss.storage)
	return &emptypb.Empty{}, nil
}

func (kvss *KeyValueStorageService) Gossip(ctx gorums.ServerCtx, request *kvsprotos.GossipMessage) (response *emptypb.Empty, err error) {
	kvss.logger.Printf("received gossip %v", request)
	kvss.mut.Lock()
	updated := kvss.storage.WriteValue(request.GetNodeID(), request.GetKey(), request.GetValue(), request.GetTimestamp())
	kvss.mut.Unlock()
	if updated {
		go func() {
			kvss.sendGossip(request.NodeID, request.GetKey())
		}()
	}
	return &emptypb.Empty{}, nil
}

type TimestampedValue struct {
	Value     int64
	Timestamp int64
}

func (tsv TimestampedValue) String() string {
	return fmt.Sprintf("(Value: %d, ts: %d)", tsv.Value, tsv.Timestamp)
}

// KeyValueStorage
// Stores key value pairs for the node and it's neighbours
//
// NodeAddress -> Key -> Value
type KeyValueStorage map[string]map[int64]TimestampedValue

func NewKeyValueStorage() *KeyValueStorage {
	return &KeyValueStorage{}
}

// ReadValue reads and combines the stored values from all neighbouring nodes
func (kvs KeyValueStorage) ReadValue(key int64) (int64, error) {
	agg := int64(0)
	found := false
	for _, values := range kvs {
		if value, exists := values[key]; exists {
			found = true
			agg += value.Value
		}
	}
	if found {
		return agg, nil
	}
	return 0, fmt.Errorf("keyvaluestorage does not contain key %v", key)
}

// ReadValueFromNode reads the stored value from a single node
func (kvs KeyValueStorage) ReadValueFromNode(nodeAddr string, key int64) (int64, error) {
	if nodeValues, containsNode := kvs[nodeAddr]; containsNode {
		if value, containsKey := nodeValues[key]; containsKey {
			return value.Value, nil
		} else {
			return 0, fmt.Errorf("keyvaluestorage does not contain key %v", key)
		}
	} else {
		return 0, fmt.Errorf("keyvaluestorage does not contain address %v", nodeAddr)
	}
}

// ReadValueExceptNode reads and combines the stored values from all neighbouring nodes except the input ID
//
// Will return (0, nil) and not (0, err) if only the excluded node has a value for the input key
func (kvs KeyValueStorage) ReadValueExceptNode(exceptID string, key int64) (int64, error) {
	agg := int64(0)
	found := false
	for nodeID, values := range kvs {
		if value, exists := values[key]; exists {
			found = true
			if exceptID != nodeID {
				agg += value.Value
			}
		}
	}
	if found {
		return agg, nil
	}
	return 0, fmt.Errorf("keyvaluestorage does not contain key %v", key)
}

// WriteValue writes the input value to storage
//
// Creates a new nodeID key if it does not exist.
// Overwrites existing key value pairs for a given nodeID.
// Returns true if the new value has a newer timestamp than last
func (kvs KeyValueStorage) WriteValue(nodeID string, key int64, value int64, timestamp int64) bool {
	if _, containsNode := kvs[nodeID]; !containsNode {
		kvs[nodeID] = make(map[int64]TimestampedValue)
	}
	if storedValue, exists := kvs[nodeID][key]; exists {
		if storedValue.Timestamp >= timestamp {
			return false
		}
	}
	kvs[nodeID][key] = TimestampedValue{Value: value, Timestamp: timestamp}
	return true
}
