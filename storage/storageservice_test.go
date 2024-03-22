package storage_test

import (
	"context"
	"testing"
	"time"

	"github.com/vidarandrebo/oncetree/storage/sqspec"

	"github.com/vidarandrebo/oncetree"
	"github.com/vidarandrebo/oncetree/consts"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/relab/gorums"
	"github.com/stretchr/testify/assert"
	kvsprotos "github.com/vidarandrebo/oncetree/protos/keyvaluestorage"
	"google.golang.org/grpc"
)

var (
	nodeIDs   = []string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9"}
	nodeAddrs = []string{
		":9080",
		":9081",
		":9082",
		":9083",
		":9084",
		":9085",
		":9086",
		":9087",
		":9088",
		":9089",
	}
	nodeMap = map[string]string{
		"0": ":9080",
		"1": ":9081",
		"2": ":9082",
		"3": ":9083",
		"4": ":9084",
		"5": ":9085",
		"6": ":9086",
		"7": ":9087",
		"8": ":9088",
		"9": ":9089",
	}
)

// TestStorageService_Write tests writing the same value to all nodes, and checking that the values has propagated to all nodes.
func TestStorageService_shareAll(t *testing.T) {
	testNodes, wg := oncetree.StartTestNodes()
	time.Sleep(consts.GorumsDialTimeout)
	mgr, cfg := createKeyValueStorageConfig()

	for _, node := range cfg.Nodes() {
		_, writeErr := node.Write(context.Background(), &kvsprotos.WriteRequest{
			Key:   20,
			Value: 10,
		})
		assert.Nil(t, writeErr)
	}
	newNode, newWg := oncetree.StartTestNode()
	time.Sleep(consts.RPCContextTimeout)

	// add the new node to config
	cfg, err := mgr.NewConfiguration(
		&sqspec.QSpec{NumNodes: len(nodeAddrs) + 1},
		cfg.WithNewNodes(gorums.WithNodeList([]string{":9090"})),
	)
	assert.Nil(t, err)

	responses, readErr := cfg.ReadAll(context.Background(), &kvsprotos.ReadRequest{
		Key: 20,
	})

	// should be as many responses as number of nodes
	assert.Equal(t, len(testNodes)+1, len(responses.GetValue()))
	assert.Nil(t, readErr)
	// checks that all nodes including the new has the value 100 for key 20
	for _, response := range responses.GetValue() {
		assert.Equal(t, int64(100), response)
	}

	newNode.Stop("stopped by test")
	for _, node := range testNodes {
		node.Stop("stopped by test")
	}
	wg.Wait()
	newWg.Wait()
}

// TestStorageService_Write tests writing the same value to all nodes, and checking that the values has propagated to all nodes.
func TestStorageService_Write(t *testing.T) {
	testNodes, wg := oncetree.StartTestNodes()
	time.Sleep(consts.GorumsDialTimeout)
	_, cfg := createKeyValueStorageConfig()

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
func createKeyValueStorageConfig() (*kvsprotos.Manager, *kvsprotos.Configuration) {
	manager := kvsprotos.NewManager(
		gorums.WithDialTimeout(consts.GorumsDialTimeout),
		gorums.WithGrpcDialOptions(
			grpc.WithBlock(),
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		),
	)
	cfg, err := manager.NewConfiguration(&sqspec.QSpec{NumNodes: len(nodeAddrs)}, gorums.WithNodeList(nodeAddrs))
	if err != nil {
		panic("failed to create cfg")
	}
	return manager, cfg
}

func BenchmarkStorageService_Write(t *testing.B) {
	testNodes, wg := oncetree.StartTestNodes()
	// time.Sleep(consts.GorumsDialTimeout)
	_, cfg := createKeyValueStorageConfig()
	node := cfg.Nodes()[0]

	for i := 0; i < t.N; i++ {
		_, err := node.Write(context.Background(), &kvsprotos.WriteRequest{
			Key:   int64(i),
			Value: int64(i),
		})
		if err != nil {
			panic(err)
		}
	}
	// time.Sleep(consts.RPCContextTimeout)

	for _, node := range testNodes {
		node.Stop("stopped by test")
	}
	wg.Wait()
}
