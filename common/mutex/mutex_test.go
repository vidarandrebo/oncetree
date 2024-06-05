package mutex

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRWMutex_Lock(t *testing.T) {
	mut := New[uint32](32)
	value := mut.Lock()
	assert.Equal(t, uint32(32), *value)
}

// TestRWMutex_Unlock checks that the lock is released properly and the reference is destroyed while preserving the value
func TestRWMutex_Unlock(t *testing.T) {
	mut := New[uint32](32)
	value := mut.Lock()
	mut.Unlock(&value)
	assert.Nil(t, value)

	reloadedValue := mut.Lock()
	assert.Equal(t, uint32(32), *reloadedValue)
}
