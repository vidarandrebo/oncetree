package oncetree_test

import (
	"context"
	"testing"
	"time"

	"google.golang.org/grpc/credentials/insecure"

	"github.com/relab/gorums"
	"github.com/stretchr/testify/assert"
	"github.com/vidarandrebo/oncetree"
	kvsprotos "github.com/vidarandrebo/oncetree/protos/keyvaluestorage"
	"google.golang.org/grpc"
)

var keyValueStorage = oncetree.KeyValueStorage{
	"addr1": {
		1: oncetree.TimestampedValue{Value: 12, Timestamp: 3},
		2: oncetree.TimestampedValue{Value: 13, Timestamp: 3},
		3: oncetree.TimestampedValue{Value: 14, Timestamp: 3},
		5: oncetree.TimestampedValue{Value: 40, Timestamp: 3},
	},
	"addr2": {
		1: oncetree.TimestampedValue{Value: 15, Timestamp: 3},
		2: oncetree.TimestampedValue{Value: 16, Timestamp: 3},
		3: oncetree.TimestampedValue{Value: 0, Timestamp: 3},
		4: oncetree.TimestampedValue{Value: 77, Timestamp: 3},
	},
	"addr3": {
		1: oncetree.TimestampedValue{Value: 44, Timestamp: 3},
		2: oncetree.TimestampedValue{Value: 88, Timestamp: 3},
		4: oncetree.TimestampedValue{Value: -77, Timestamp: 3},
	},
	"addr4": {
		1: oncetree.TimestampedValue{Value: 15, Timestamp: 3},
		2: oncetree.TimestampedValue{Value: 16, Timestamp: 3},
		3: oncetree.TimestampedValue{Value: 0, Timestamp: 3},
		4: oncetree.TimestampedValue{Value: 55, Timestamp: 3},
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
	valueChanged := keyValueStorage.WriteValue(testAddr, testKey, testValue, 4)
	assert.Equal(t, keyValueStorage[testAddr][testKey], oncetree.TimestampedValue{Value: testValue, Timestamp: 4})
	assert.True(t, valueChanged)
}

func TestKeyValueStorage_WriteValue_OverWrite(t *testing.T) {
	testAddr := "addr1"
	testKey := int64(2)
	testValue := int64(15)
	valueChanged := keyValueStorage.WriteValue(testAddr, testKey, testValue, 4)
	assert.Equal(t, keyValueStorage[testAddr][testKey], oncetree.TimestampedValue{Value: testValue, Timestamp: 4})
	assert.True(t, valueChanged)
}

// Value should not change if ts is lower than the stored value
func TestKeyValueStorage_WriteValue_NoChange(t *testing.T) {
	testAddr := "addr2"
	testKey := int64(2)
	testValue := int64(99)
	testTs := int64(2)
	valueChanged := keyValueStorage.WriteValue(testAddr, testKey, testValue, testTs)
	assert.Equal(t, keyValueStorage[testAddr][testKey], oncetree.TimestampedValue{Value: 16, Timestamp: 3})
	assert.False(t, valueChanged)
}

// TestKeyValueStorageService_Write tests writing the same value to all nodes, and checking that the values has propagated to all nodes.
func TestKeyValueStorageService_Write(t *testing.T) {
	testNodes, wg := StartTestNodes()
	time.Sleep(1 * time.Second)
	cfg := createKeyValueStorageConfig()

	for _, node := range cfg.Nodes() {
		_, writeErr := node.Write(context.Background(), &kvsprotos.WriteRequest{
			Key:   20,
			Value: 10,
		})
		assert.Nil(t, writeErr)
	}
	time.Sleep(1 * time.Second)

	responses, readErr := cfg.ReadAll(context.Background(), &kvsprotos.ReadRequest{
		Key: 20,
	})

	// should be as many responses as number of nodes
	assert.Equal(t, len(testNodes), len(responses.GetValue()))
	assert.Nil(t, readErr)
	for _, response := range responses.GetValue() {
		assert.Equal(t, int64(100), response)
	}

	for _, node := range testNodes {
		node.Stop("stopped by test")
	}
	wg.Wait()
}

// createKeyValueStorageConfig creates a new manager and returns an initialized configuration to use with the KeyValueStorageService
func createKeyValueStorageConfig() *kvsprotos.Configuration {
	manager := kvsprotos.NewManager(
		gorums.WithDialTimeout(1*time.Second),
		gorums.WithGrpcDialOptions(
			grpc.WithBlock(),
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		),
	)
	cfg, err := manager.NewConfiguration(&oncetree.QSpec{NumNodes: len(nodeAddrs)}, gorums.WithNodeList(nodeAddrs))
	if err != nil {
		panic("failed to create cfg")
	}
	return cfg
}
