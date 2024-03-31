package hashset

import (
	"math/rand"
	"testing"
)

func BenchmarkConcurrentHashSet_String(b *testing.B) {
	set := NewHashSet[int]()
	for len(set) < 1000 {
		set.Add(rand.Int())
	}
	for i := 0; i < b.N; i++ {
		_ = set.String()
	}
}

func BenchmarkConcurrentHashSet_Values(b *testing.B) {
	set := NewHashSet[int]()
	for len(set) < 1000 {
		set.Add(rand.Int())
	}
	for i := 0; i < b.N; i++ {
		set.Values()
	}
}
