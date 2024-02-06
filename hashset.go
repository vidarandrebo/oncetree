package oncetree

type (
	HashSet[T comparable] map[T]struct{}
)

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
