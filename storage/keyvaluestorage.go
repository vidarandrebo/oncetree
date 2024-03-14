package storage

import (
	"fmt"
	"sync"

	"github.com/vidarandrebo/oncetree/concurrent/hashset"
)

// KeyValueStorage
// Stores key value pairs for the node and it's neighbours
//
// NodeAddress -> Key -> Value
type KeyValueStorage struct {
	data map[string]map[int64]TimestampedValue
	mut  sync.RWMutex
}

func NewKeyValueStorage() *KeyValueStorage {
	return &KeyValueStorage{
		data: make(map[string]map[int64]TimestampedValue),
	}
}

// ReadValue reads and combines the stored values from all neighbouring nodes
func (kvs *KeyValueStorage) ReadValue(key int64) (int64, error) {
	agg := int64(0)
	found := false
	for _, values := range kvs.data {
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
func (kvs *KeyValueStorage) ReadValueFromNode(nodeID string, key int64) (int64, error) {
	if nodeValues, containsNode := kvs.data[nodeID]; containsNode {
		if value, containsKey := nodeValues[key]; containsKey {
			return value.Value, nil
		} else {
			return 0, fmt.Errorf("keyvaluestorage does not contain key %v", key)
		}
	} else {
		return 0, fmt.Errorf("keyvaluestorage does not contain address %v", nodeID)
	}
}

// ReadValueExceptNode reads and combines the stored values from all neighbouring nodes except the input ID
//
// Will return (0, nil) and not (0, err) if only the excluded node has a value for the input key
func (kvs *KeyValueStorage) ReadValueExceptNode(exceptID string, key int64) (int64, error) {
	agg := int64(0)
	found := false
	for nodeID, values := range kvs.data {
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
func (kvs *KeyValueStorage) WriteValue(nodeID string, key int64, value int64, timestamp int64) bool {
	if _, containsNode := kvs.data[nodeID]; !containsNode {
		kvs.data[nodeID] = make(map[int64]TimestampedValue)
	}
	if storedValue, exists := kvs.data[nodeID][key]; exists {
		if storedValue.Timestamp >= timestamp {
			return false
		}
	}
	kvs.data[nodeID][key] = TimestampedValue{Value: value, Timestamp: timestamp}
	return true
}

func (kvs *KeyValueStorage) Keys() hashset.HashSet[int64] {
	keys := hashset.NewHashSet[int64]()
	for _, nodeValues := range kvs.data {
		for key := range nodeValues {
			keys.Add(key)
		}
	}
	return keys
}
