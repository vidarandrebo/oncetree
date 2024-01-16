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
	data := KeyValueStorage{}
	return &data
}

// ReadValue reads and combines the values from all neighbouring nodes
func (kvs KeyValueStorage) ReadValue(key int64) (int64, error) {
	return 0, fmt.Errorf("not implemented")
}

// ReadValueFromNode reads the value from a single node
func (kvs KeyValueStorage) ReadValueFromNode(addr string, key int64) (int64, error) {
	if nodeValues, containsNode := kvs[addr]; containsNode {
		if value, containsKey := nodeValues[key]; containsKey {
			return value, nil
		} else {
			return 0, fmt.Errorf("keyvaluestorage does not contain key %v", key)
		}
	} else {
		return 0, fmt.Errorf("keyvaluestorage does not contain address %v", addr)
	}
}

// ReadValueExceptNode reads and combines the values from all neighbouring nodes except the input addr
func (kvs KeyValueStorage) ReadValueExceptNode(addr string, key int64) (int64, error) {
	return 0, fmt.Errorf("not implemented")
}

func (kvs KeyValueStorage) WriteValue(addr string, key int64, value int64) error {
	return fmt.Errorf("not implemented")
}
