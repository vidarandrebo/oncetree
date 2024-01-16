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
		5: 40,
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

func TestKeyValueStorage_ReadValueFromNode(t *testing.T) {
	value, err := keyValueStorage.ReadValueFromNode("addr1", 1)
	assert.Nil(t, err)
	assert.Equal(t, int64(12), value)
}
func TestKeyValueStorage_ReadValueFromNode_NoAddr(t *testing.T) {
	value, err := keyValueStorage.ReadValueFromNode("addr99", 1)
	assert.Equal(t, int64(0), value)
	assert.NotNil(t, err)
}
func TestKeyValueStorage_ReadValueFromNode_NoKey(t *testing.T) {
	value, err := keyValueStorage.ReadValueFromNode("addr1", 99)
	assert.Equal(t, int64(0), value)
	assert.NotNil(t, err)
}

func TestKeyValueStorage_ReadValue(t *testing.T) {
	value, err := keyValueStorage.ReadValue(4)
	assert.Nil(t, err)
	assert.Equal(t, int64(55), value)

	value, err = keyValueStorage.ReadValue(1)
	assert.Nil(t, err)
	assert.Equal(t, int64(86), value)
}

func TestKeyValueStorage_ReadValue_NoKey(t *testing.T) {
	value, err := keyValueStorage.ReadValue(99)
	assert.NotNil(t, err)
	assert.Equal(t, int64(0), value)
}
func TestKeyValueStorage_ReadValueExceptNode(t *testing.T) {
	value, err := keyValueStorage.ReadValueExceptNode("addr1", 4)
	assert.Nil(t, err)
	assert.Equal(t, int64(55), value)

	value, err = keyValueStorage.ReadValueExceptNode("addr4", 1)
	assert.Nil(t, err)
	assert.Equal(t, int64(71), value)
}

func TestKeyValueStorage_ReadValueExceptNode_NoKey(t *testing.T) {
	value, err := keyValueStorage.ReadValueExceptNode("addr1", 99)
	assert.NotNil(t, err)
	assert.Equal(t, int64(0), value)
}

func TestKeyValueStorage_ReadValueExceptNode_FoundZero(t *testing.T) {
	value, err := keyValueStorage.ReadValueExceptNode("addr1", 5)
	assert.Nil(t, err)
	assert.Equal(t, int64(0), value)
}

func TestKeyValueStorage_WriteValue_NoAddr(t *testing.T) {
	testAddr := "addr55"
	testKey := int64(1)
	testValue := int64(10)
	valueChanged := keyValueStorage.WriteValue(testAddr, testKey, testValue)
	assert.Equal(t, keyValueStorage[testAddr][testKey], testValue)
	assert.True(t, valueChanged)
}
func TestKeyValueStorage_WriteValue_OverWrite(t *testing.T) {
	testAddr := "addr1"
	testKey := int64(2)
	testValue := int64(15)
	valueChanged := keyValueStorage.WriteValue(testAddr, testKey, testValue)
	assert.Equal(t, keyValueStorage[testAddr][testKey], testValue)
	assert.True(t, valueChanged)
}
func TestKeyValueStorage_WriteValue_NoChange(t *testing.T) {
	testAddr := "addr2"
	testKey := int64(2)
	testValue := int64(16)
	valueChanged := keyValueStorage.WriteValue(testAddr, testKey, testValue)
	assert.Equal(t, keyValueStorage[testAddr][testKey], testValue)
	assert.False(t, valueChanged)
}
