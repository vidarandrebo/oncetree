package oncetree

import (
	"fmt"
)

// KeyValueStorage
// Stores key value pairs for the node and it's neighbours
//
// NodeAddress -> Key -> Value
type KeyValueStorage map[string]map[int64]int64

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
			agg += value
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
			return value, nil
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
				agg += value
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
// Returns true if value has changed.
func (kvs KeyValueStorage) WriteValue(nodeID string, key int64, value int64) bool {
	if _, containsNode := kvs[nodeID]; !containsNode {
		kvs[nodeID] = make(map[int64]int64)
	}
	if storedValue, exists := kvs[nodeID][key]; exists {
		if storedValue == value {
			return false
		}
	}
	kvs[nodeID][key] = value
	return true
}
