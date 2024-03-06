package maps

import "sync"

type ConcurrentMap[TKey comparable, TVal any] struct {
	rwLock sync.RWMutex
	data   map[TKey]TVal
}

func NewConcurrentMap[TKey comparable, TVal any]() *ConcurrentMap[TKey, TVal] {
	return &ConcurrentMap[TKey, TVal]{
		data: make(map[TKey]TVal),
	}
}

func (cm *ConcurrentMap[TKey, TVal]) Set(key TKey, value TVal) {
	cm.rwLock.Lock()
	defer cm.rwLock.Unlock()
	cm.data[key] = value
}

func (cm *ConcurrentMap[TKey, TVal]) Get(key TKey) (TVal, bool) {
	cm.rwLock.RLock()
	defer cm.rwLock.RUnlock()
	value, ok := cm.data[key]
	return value, ok
}

func (cm *ConcurrentMap[TKey, TVal]) Delete(key TKey) {
	cm.rwLock.Lock()
	defer cm.rwLock.Unlock()
	delete(cm.data, key)
}

func (cm *ConcurrentMap[TKey, TVal]) Keys() []TKey {
	keys := make([]TKey, 0)
	cm.rwLock.RLock()
	defer cm.rwLock.RUnlock()
	for key := range cm.data {
		keys = append(keys, key)
	}
	return keys
}

func (cm *ConcurrentMap[TKey, TVal]) Values() []TVal {
	values := make([]TVal, 0)
	cm.rwLock.RLock()
	defer cm.rwLock.RUnlock()
	for _, value := range cm.data {
		values = append(values, value)
	}
	return values
}

type KeyValuePair[TKey any, TVal any] struct {
	Key   TKey
	Value TVal
}

func (cm *ConcurrentMap[TKey, TVal]) Entries() []KeyValuePair[TKey, TVal] {
	entries := make([]KeyValuePair[TKey, TVal], 0)
	cm.rwLock.RLock()
	defer cm.rwLock.RUnlock()
	for key, value := range cm.data {
		entries = append(entries, KeyValuePair[TKey, TVal]{Key: key, Value: value})
	}
	return entries
}

func (cm *ConcurrentMap[TKey, TVal]) Contains(key TKey) bool {
	cm.rwLock.RLock()
	defer cm.rwLock.RUnlock()
	_, ok := cm.data[key]
	return ok
}

func (cm *ConcurrentMap[TKey, TVal]) Clear() {
	cm.rwLock.Lock()
	defer cm.rwLock.Unlock()
	for key := range cm.data {
		delete(cm.data, key)
	}
}

type ConcurrentIntegerMap[TKey comparable] struct {
	ConcurrentMap[TKey, int]
}

func NewConcurrentIntegerMap[TKey comparable]() *ConcurrentIntegerMap[TKey] {
	return &ConcurrentIntegerMap[TKey]{
		ConcurrentMap[TKey, int]{
			data: make(map[TKey]int),
		},
	}
}

func (cim *ConcurrentIntegerMap[TKey]) Increment(key TKey, amount int) {
	cim.rwLock.Lock()
	defer cim.rwLock.Unlock()
	_, ok := cim.data[key]
	if !ok {
		cim.data[key] = amount
		return
	}
	cim.data[key] += amount
}
