package oncetree

import "sync"

type HashSet[T comparable] map[T]struct{}

func NewHashSet[T comparable]() HashSet[T] {
	return make(HashSet[T])
}

func (s HashSet[T]) Add(key T) {
	s[key] = struct{}{}
}

func (s HashSet[T]) Contains(key T) bool {
	_, ok := s[key]
	return ok
}

func (s HashSet[T]) Remove(key T) {
	delete(s, key)
}
func (s HashSet[T]) Values() []T {
	values := make([]T, 0)
	for value := range s {
		values = append(values, value)
	}
	return values
}
func (s HashSet[T]) Clear() {
	for k := range s {
		delete(s, k)
	}
}

func (s HashSet[T]) Intersection(b HashSet[T]) HashSet[T] {
	result := NewHashSet[T]()
	for key := range s {
		if b.Contains(key) {
			result.Add(key)
		}
	}
	return result
}

func (s HashSet[T]) Union(b HashSet[T]) HashSet[T] {
	panic("not implemented")
}

type ConcurrentHashSet[T comparable] struct {
	hashSet HashSet[T]
	mut     sync.RWMutex
}

func NewConcurrentHashSet[T comparable]() *ConcurrentHashSet[T] {
	return &ConcurrentHashSet[T]{
		hashSet: make(HashSet[T]),
	}
}
func (s *ConcurrentHashSet[T]) Len() int {
	s.mut.RLock()
	defer s.mut.RUnlock()
	return len(s.hashSet)
}

func (s *ConcurrentHashSet[T]) Values() []T {
	s.mut.RLock()
	defer s.mut.RUnlock()
	return s.hashSet.Values()
}
func (s *ConcurrentHashSet[T]) Add(key T) {
	s.mut.Lock()
	defer s.mut.Unlock()
	s.hashSet.Add(key)
}

func (s *ConcurrentHashSet[T]) Clear() {
	s.mut.Lock()
	defer s.mut.Unlock()
	s.hashSet.Clear()
}

func (s *ConcurrentHashSet[T]) Contains(key T) bool {
	s.mut.RLock()
	defer s.mut.RUnlock()
	return s.hashSet.Contains(key)
}

func (s *ConcurrentHashSet[T]) Remove(key T) {
	s.mut.Lock()
	defer s.mut.Unlock()
	s.hashSet.Remove(key)
}

func (s *ConcurrentHashSet[T]) Intersection(b *ConcurrentHashSet[T]) *ConcurrentHashSet[T] {
	s.mut.RLock()
	b.mut.RLock()
	defer s.mut.RUnlock()
	defer b.mut.RUnlock()
	return &ConcurrentHashSet[T]{
		hashSet: s.hashSet.Intersection(b.hashSet),
	}

}

func (s *ConcurrentHashSet[T]) Union(b *ConcurrentHashSet[T]) ConcurrentHashSet[T] {
	panic("not implemented")
}
