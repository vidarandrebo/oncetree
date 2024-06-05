package mutex

import "sync"

// RWMutex encapsulates a value of type T and a read-write mutex
type RWMutex[T any] struct {
	mut   sync.RWMutex
	value T
}

func New[T any](value T) *RWMutex[T] {
	return &RWMutex[T]{
		value: value,
	}
}

// Lock locks the value for reading and writing, and returns a pointer to the value
func (rwm *RWMutex[T]) Lock() *T {
	rwm.mut.Lock()
	return &rwm.value
}

// RLock locks the value for reading and returns a pointer to the value
func (rwm *RWMutex[T]) RLock() *T {
	rwm.mut.RLock()
	return &rwm.value
}

// RUnlock releases the lock and nulls the reference provided in the argument
func (rwm *RWMutex[T]) RUnlock(pt **T) {
	*pt = nil
	rwm.mut.RUnlock()
}

// Unlock releases the lock and nulls the reference provided in the argument
func (rwm *RWMutex[T]) Unlock(pt **T) {
	*pt = nil
	rwm.mut.Unlock()
}
