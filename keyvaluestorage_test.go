package oncetree

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

var keyValueStorage = KeyValueStorage{
	"addr1": {
		1: 12,
		2: 13,
		3: 14,
	},
	"addr2": {
		1: 15,
		2: 16,
		3: 0,
		4: 77,
	},
	"addr3": {
		1: 44,
		2: 88,
		4: -77,
	},
	"addr4": {
		1: 15,
		2: 16,
		3: 0,
		4: 55,
	},
}

func TestReadValueFromNode(t *testing.T) {
	value, err := keyValueStorage.ReadValueFromNode("addr1", 1)
	assert.Nil(t, err)
	assert.Equal(t, int64(12), value)
}
func TestReadValueFromNodeNoAddr(t *testing.T) {
	_, err := keyValueStorage.ReadValueFromNode("addr99", 1)
	assert.NotNil(t, err)
}
func TestReadValueFromNodeNoKey(t *testing.T) {
	_, err := keyValueStorage.ReadValueFromNode("addr1", 99)
	assert.NotNil(t, err)
}
