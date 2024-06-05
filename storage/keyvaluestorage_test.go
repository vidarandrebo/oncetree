package storage

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var keyValueStorage = KeyValueStorage{
	data: map[string]map[int64]TimestampedValue{
		"ID1": {
			1: TimestampedValue{Value: 12, Timestamp: 3},
			2: TimestampedValue{Value: 13, Timestamp: 3},
			3: TimestampedValue{Value: 14, Timestamp: 3},
			5: TimestampedValue{Value: 40, Timestamp: 3},
		},
		"ID2": {
			1: TimestampedValue{Value: 15, Timestamp: 3},
			2: TimestampedValue{Value: 16, Timestamp: 3},
			3: TimestampedValue{Value: 0, Timestamp: 3},
			4: TimestampedValue{Value: 77, Timestamp: 3},
		},
		"ID3": {
			1: TimestampedValue{Value: 44, Timestamp: 3},
			2: TimestampedValue{Value: 88, Timestamp: 3},
			4: TimestampedValue{Value: -77, Timestamp: 3},
		},
		"ID4": {
			1: TimestampedValue{Value: 15, Timestamp: 3},
			2: TimestampedValue{Value: 16, Timestamp: 3},
			3: TimestampedValue{Value: 0, Timestamp: 3},
			4: TimestampedValue{Value: 55, Timestamp: 3},
		},
	},
}

func TestKeyValueStorage_ReadValueFromNode(t *testing.T) {
	value, err := keyValueStorage.ReadValueFromNode(1, "ID1")
	assert.Nil(t, err)
	assert.Equal(t, int64(12), value.Value)
}

func TestKeyValueStorage_ReadValueFromNode_NoAddr(t *testing.T) {
	value, err := keyValueStorage.ReadValueFromNode(1, "ID99")
	assert.NotNil(t, err)
	assert.Equal(t, int64(0), value.Value)
}

func TestKeyValueStorage_ReadValueFromNode_NoKey(t *testing.T) {
	value, err := keyValueStorage.ReadValueFromNode(99, "ID1")
	assert.NotNil(t, err)
	assert.Equal(t, int64(0), value.Value)
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
	value, err := keyValueStorage.GossipValue(4, "ID1")
	assert.Nil(t, err)
	assert.Equal(t, int64(55), value)

	value, err = keyValueStorage.GossipValue(1, "ID4")
	assert.Nil(t, err)
	assert.Equal(t, int64(71), value)
}

func TestKeyValueStorage_ReadValueExceptNode_NoKey(t *testing.T) {
	value, err := keyValueStorage.GossipValue(99, "ID1")
	assert.NotNil(t, err)
	assert.Equal(t, int64(0), value)
}

func TestKeyValueStorage_ReadValueExceptNode_FoundZero(t *testing.T) {
	value, err := keyValueStorage.GossipValue(5, "ID1")
	assert.Nil(t, err)
	assert.Equal(t, int64(0), value)
}

func TestKeyValueStorage_WriteValue_NoAddr(t *testing.T) {
	testID := "ID55"
	testKey := int64(1)
	testValue := int64(10)
	testTs := int64(4)
	valueChanged := keyValueStorage.WriteValue(testKey, testValue, testTs, testID)
	assert.Equal(t, TimestampedValue{Value: testValue, Timestamp: 4}, keyValueStorage.data[testID][testKey])
	assert.True(t, valueChanged)
}

func TestKeyValueStorage_WriteValue_OverWrite(t *testing.T) {
	testID := "ID1"
	testKey := int64(2)
	testValue := int64(15)
	testTs := int64(4)
	valueChanged := keyValueStorage.WriteValue(testKey, testValue, testTs, testID)
	assert.Equal(t, TimestampedValue{Value: testValue, Timestamp: 4}, keyValueStorage.data[testID][testKey])
	assert.True(t, valueChanged)
}

// Value should not change if ts is lower than the stored value
func TestKeyValueStorage_WriteValue_NoChange(t *testing.T) {
	testID := "ID2"
	testKey := int64(2)
	testValue := int64(99)
	testTs := int64(2)
	valueChanged := keyValueStorage.WriteValue(testKey, testValue, testTs, testID)
	assert.Equal(t, keyValueStorage.data[testID][testKey], TimestampedValue{Value: 16, Timestamp: 3})
	assert.False(t, valueChanged)
}
