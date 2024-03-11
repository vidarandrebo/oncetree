package storage_test

import (
	"context"
	"testing"
	"time"

	"github.com/vidarandrebo/oncetree"
	"github.com/vidarandrebo/oncetree/consts"
	"github.com/vidarandrebo/oncetree/storage"

	"google.golang.org/grpc/credentials/insecure"

	"github.com/relab/gorums"
	"github.com/stretchr/testify/assert"
	kvsprotos "github.com/vidarandrebo/oncetree/protos/keyvaluestorage"
	"google.golang.org/grpc"
)

// TestStorageService_Write tests writing the same value to all nodes, and checking that the values has propagated to all nodes.
func TestStorageService_Write(t *testing.T) {
	testNodes, wg := oncetree.StartTestNodes()
	time.Sleep(consts.GorumsDialTimeout)
	cfg := createKeyValueStorageConfig()

	for _, node := range cfg.Nodes() {
		_, writeErr := node.Write(context.Background(), &kvsprotos.WriteRequest{
			Key:   20,
			Value: 10,
		})
		assert.Nil(t, writeErr)
	}
	time.Sleep(consts.RPCContextTimeout)

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

// createKeyValueStorageConfig creates a new manager and returns an initialized configuration to use with the StorageService
func createKeyValueStorageConfig() *kvsprotos.Configuration {
	manager := kvsprotos.NewManager(
		gorums.WithDialTimeout(consts.GorumsDialTimeout),
		gorums.WithGrpcDialOptions(
			grpc.WithBlock(),
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		),
	)
	cfg, err := manager.NewConfiguration(&storage.QSpec{NumNodes: len(nodeAddrs)}, gorums.WithNodeList(nodeAddrs))
	if err != nil {
		panic("failed to create cfg")
	}
	return cfg
}
