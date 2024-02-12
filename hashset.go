package oncetree

import "sync"

type ConcurrentHashSet[T comparable] struct {
	set map[T]struct{}
	mut sync.RWMutex
}

func NewHashSet[T comparable]() *ConcurrentHashSet[T] {
	return &ConcurrentHashSet[T]{
		set: make(map[T]struct{}),
	}
}
func (s *ConcurrentHashSet[T]) Len() int {
	s.mut.RLock()
	defer s.mut.RUnlock()
	return len(s.set)
}

func (s *ConcurrentHashSet[T]) Values() []T {
	s.mut.RLock()
	defer s.mut.RUnlock()
	values := make([]T, 0)
	for value := range s.set {
		values = append(values, value)
	}
	return values
}
func (s *ConcurrentHashSet[T]) Add(key T) {
	s.mut.Lock()
	defer s.mut.Unlock()
	s.set[key] = struct{}{}
}

func (s *ConcurrentHashSet[T]) Contains(key T) bool {
	s.mut.RLock()
	defer s.mut.RUnlock()
	_, ok := s.set[key]
	return ok
}

func (s *ConcurrentHashSet[T]) Remove(key T) {
	s.mut.Lock()
	defer s.mut.Unlock()
	delete(s.set, key)
}

func (s *ConcurrentHashSet[T]) Intersection(b *ConcurrentHashSet[T]) *ConcurrentHashSet[T] {
	s.mut.RLock()
	defer s.mut.RUnlock()
	result := NewHashSet[T]()
	for key := range s.set {
		if b.Contains(key) {
			result.Add(key)
		}
	}
	return result
}

func (s *ConcurrentHashSet[T]) Union(b *ConcurrentHashSet[T]) ConcurrentHashSet[T] {
	panic("not implemented")
}
