package storage

import (
	"fmt"
	"log"
	"sync"

	"github.com/vidarandrebo/oncetree/common/hashset"
)

// KeyValueStorage
// Stores key value pairs for the node and it's neighbours
//
// NodeAddress -> Key -> Value
type KeyValueStorage struct {
	data             map[string]map[int64]TimestampedValue
	local            map[string]map[int64]TimestampedValue
	readExclusions   map[int64]hashset.HashSet[string]            // key -> nodeID
	gossipExclusions map[int64]map[string]hashset.HashSet[string] // key -> dest -> nodeID
	mut              sync.RWMutex
}

func NewKeyValueStorage() *KeyValueStorage {
	return &KeyValueStorage{
		data:             make(map[string]map[int64]TimestampedValue),
		local:            make(map[string]map[int64]TimestampedValue),
		readExclusions:   make(map[int64]hashset.HashSet[string]),
		gossipExclusions: make(map[int64]map[string]hashset.HashSet[string]),
	}
}

// ReadValue reads and combines the stored values from all neighbouring nodes
func (kvs *KeyValueStorage) ReadValue(key int64) (int64, error) {
	kvs.mut.RLock()
	defer kvs.mut.RUnlock()
	agg := int64(0)
	found := false
	for nodeID, values := range kvs.data {
		if value, exists := values[key]; exists {
			found = true
			// we skip the key for that node if it is excluded
			if excludedNodes, ok := kvs.readExclusions[key]; ok && excludedNodes.Contains(nodeID) {
				continue
			}
			agg += value.Value
		}
	}
	if found {
		return agg, nil
	}
	return 0, fmt.Errorf("keyvaluestorage does not contain key %v", key)
}

// ReadLocalValue reads the local value of key from the input nodeID
func (kvs *KeyValueStorage) ReadLocalValue(key int64, nodeID string) (TimestampedValue, bool) {
	kvs.mut.RLock()
	defer kvs.mut.RUnlock()
	value, ok := kvs.local[nodeID][key]
	if !ok {
		return TimestampedValue{
			Value:     0,
			Timestamp: 0,
		}, false
	}
	return value, true
}

// ReadValueFromNode reads the stored value from a single node
func (kvs *KeyValueStorage) ReadValueFromNode(key int64, nodeID string) (TimestampedValue, error) {
	kvs.mut.RLock()
	defer kvs.mut.RUnlock()
	if nodeValues, containsNode := kvs.data[nodeID]; containsNode {
		if value, containsKey := nodeValues[key]; containsKey {
			return value, nil
		} else {
			return TimestampedValue{}, fmt.Errorf("keyvaluestorage does not contain key %v", key)
		}
	} else {
		return TimestampedValue{}, fmt.Errorf("keyvaluestorage does not contain id %v", nodeID)
	}
}

// GossipValue reads and combines the stored values from all neighbouring nodes except the input ID
//
// Will return (0, nil) and not (0, err) if only the excluded node has a value for the input key
func (kvs *KeyValueStorage) GossipValue(key int64, destId string) (int64, error) {
	kvs.mut.RLock()
	defer kvs.mut.RUnlock()
	agg := int64(0)
	found := false
	for nodeID, values := range kvs.data {
		if value, exists := values[key]; exists {
			found = true

			// exclusions contain the key
			if destinations, hasExcludeForKey := kvs.gossipExclusions[key]; hasExcludeForKey {
				// exclusions contain the destination of gossip
				if excludes, hasExcludeForDestination := destinations[destId]; hasExcludeForDestination && excludes.Contains(nodeID) {
					continue
				}
			}
			if destId != nodeID {
				agg += value.Value
			}
		}
	}
	if found {
		return agg, nil
	}
	return 0, fmt.Errorf("keyvaluestorage does not contain key %v", key)
}

// GossipValues returns a map nodeID -> value based on the input key
func (kvs *KeyValueStorage) GossipValues(key int64, nodeIDs []string) (map[string]int64, error) {
	values := make(map[string]int64)
	for _, nodeID := range nodeIDs {
		// GossipValue grabs lock, no need to lock here
		value, err := kvs.GossipValue(key, nodeID)
		if err != nil {
			return nil, err
		}
		values[nodeID] = value
	}
	return values, nil
}

// WriteValue writes the input value to storage
//
// Creates a new nodeID key if it does not exist.
// Overwrites existing key value pairs for a given nodeID.
// Returns true if the new value has a newer timestamp than last
func (kvs *KeyValueStorage) WriteValue(key int64, value int64, timestamp int64, nodeID string) bool {
	kvs.mut.Lock()
	defer kvs.mut.Unlock()
	return kvs.writeValue(key, value, timestamp, nodeID)
}

func (kvs *KeyValueStorage) writeValue(key int64, value int64, timestamp int64, nodeID string) bool {
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

func (kvs *KeyValueStorage) WriteLocalValue(key int64, value int64, timestamp int64, nodeID string) bool {
	kvs.mut.Lock()
	defer kvs.mut.Unlock()
	return kvs.writeLocalValue(key, value, timestamp, nodeID)
}

func (kvs *KeyValueStorage) writeLocalValue(key int64, value int64, timestamp int64, nodeID string) bool {
	// timestamp == 0 is the default and a write operation has not occured at nodeID
	if timestamp == 0 {
		return false
	}
	if _, containsNode := kvs.local[nodeID]; !containsNode {
		kvs.local[nodeID] = make(map[int64]TimestampedValue)
	}
	if storedValue, exists := kvs.local[nodeID][key]; exists {
		if storedValue.Timestamp >= timestamp {
			return false
		}
	}
	kvs.local[nodeID][key] = TimestampedValue{Value: value, Timestamp: timestamp}
	return true
}

// DeleteLocal removes the local key-value pair for key with owner nodeID
// Thread safe
func (kvs *KeyValueStorage) DeleteLocal(key int64, nodeID string) {
	kvs.mut.Lock()
	defer kvs.mut.Unlock()
	kvs.deleteLocal(key, nodeID)
}

func (kvs *KeyValueStorage) deleteLocal(key int64, nodeID string) {
	delete(kvs.data[nodeID], key)
}

// DeleteAgg removes the key-value pair for key with owner nodeID
// Thread safe
func (kvs *KeyValueStorage) DeleteAgg(key int64, nodeID string) {
	kvs.mut.Lock()
	defer kvs.mut.Unlock()
	kvs.deleteAgg(key, nodeID)
}

func (kvs *KeyValueStorage) deleteAgg(key int64, nodeID string) {
	delete(kvs.local[nodeID], key)
}

// ChangeKeyValueOwnership removes the key-value pair for old-owner and updates the key-value pair for the new owner
func (kvs *KeyValueStorage) ChangeKeyValueOwnership(key int64, aggValue int64, localValue int64, ts int64, oldOwner string, newOwner string) {
	kvs.mut.Lock()
	defer kvs.mut.Unlock()
	wroteLocal := kvs.writeLocalValue(key, localValue, ts, newOwner)
	if !wroteLocal {
		log.Println("did not write local value")
	}
	wroteValue := kvs.writeValue(key, aggValue, ts, newOwner)
	if !wroteValue {
		fmt.Println("did not write value")
	}
	kvs.deleteAgg(key, oldOwner)
	kvs.deleteLocal(key, oldOwner)
}

func (kvs *KeyValueStorage) Keys() hashset.HashSet[int64] {
	kvs.mut.RLock()
	defer kvs.mut.RUnlock()
	return kvs.keys()
}

func (kvs *KeyValueStorage) keys() hashset.HashSet[int64] {
	keys := hashset.New[int64]()
	for _, nodeValues := range kvs.data {
		for key := range nodeValues {
			keys.Add(key)
		}
	}
	return keys
}

func (kvs *KeyValueStorage) AddReadExclusions(excludeIDs ...string) {
	kvs.mut.Lock()
	kvs.addReadExclusions(excludeIDs...)
	kvs.mut.Unlock()
}

func (kvs *KeyValueStorage) addReadExclusions(excludeIDs ...string) {
	for _, key := range kvs.keys().Values() {
		kvs.readExclusions[key] = hashset.New[string]()
		for _, id := range excludeIDs {
			kvs.readExclusions[key].Add(id)
		}
	}
}

func (kvs *KeyValueStorage) AddGossipExclusions(excludes map[string]hashset.HashSet[string]) {
	kvs.mut.Lock()
	kvs.addGossipExclusions(excludes)
	kvs.mut.Unlock()
}

func (kvs *KeyValueStorage) addGossipExclusions(excludes map[string]hashset.HashSet[string]) {
	for _, key := range kvs.keys().Values() {
		kvs.gossipExclusions[key] = make(map[string]hashset.HashSet[string])
		for dest, exclude := range excludes {
			kvs.gossipExclusions[key][dest] = exclude
		}
	}
}

func (kvs *KeyValueStorage) RemoveExclusionsFromKey(key int64) {
	kvs.mut.Lock()
	defer kvs.mut.Unlock()
	delete(kvs.gossipExclusions, key)
	delete(kvs.readExclusions, key)
}
