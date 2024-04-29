package storage_test

import (
	"context"
	"fmt"
	"log/slog"
	"maps"
	"os"
	"testing"
	"time"

	"github.com/vidarandrebo/oncetree/gorumsprovider"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/stretchr/testify/assert"
	"github.com/vidarandrebo/oncetree"
	"github.com/vidarandrebo/oncetree/consts"
	kvsprotos "github.com/vidarandrebo/oncetree/protos/keyvaluestorage"
	nodeprotos "github.com/vidarandrebo/oncetree/protos/node"
)

// structure same as in testing_node.go
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
	gorumsNodeMap = map[string]uint32{
		":9080": 0,
		":9081": 1,
		":9082": 2,
		":9083": 3,
		":9084": 4,
		":9085": 5,
		":9086": 6,
		":9087": 7,
		":9088": 8,
		":9089": 9,
	}
	logger = slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		AddSource:   false,
		Level:       consts.LogLevel,
		ReplaceAttr: nil,
	})).With(slog.Group("node", slog.String("module", "storage_test")))
)

func nodeConfig(provider *gorumsprovider.GorumsProvider) (*nodeprotos.Configuration, map[string]uint32) {
	nodeMap := make(map[string]uint32)
	maps.Copy(nodeMap, gorumsNodeMap)
	provider.SetNodes(nodeMap)
	cfg, _ := provider.NodeConfig()
	return cfg, nodeMap
}

func storageConfig(provider *gorumsprovider.GorumsProvider) (*kvsprotos.Configuration, map[string]uint32) {
	nodeMap := make(map[string]uint32)
	maps.Copy(nodeMap, gorumsNodeMap)
	provider.SetNodes(nodeMap)
	cfg, _ := provider.StorageConfig()
	return cfg, nodeMap
}

// TestStorageService_Write tests writing the same value to all nodes, and checking that the values has propagated to all nodes.
func TestStorageService_shareAll(t *testing.T) {
	testNodes, wg := oncetree.StartTestNodes(true)

	gorumsProvider := gorumsprovider.New(logger)
	cfg, nodeMap := storageConfig(gorumsProvider)

	for _, node := range cfg.Nodes() {
		_, writeErr := node.Write(context.Background(), &kvsprotos.WriteRequest{
			Key:   20,
			Value: 10,
		})
		assert.Nil(t, writeErr)
	}
	newNode, newWg := oncetree.StartTestNode(true)
	time.Sleep(consts.RPCContextTimeout)

	// add the new node to config
	nodeMap[":9090"] = 10
	gorumsProvider.ResetWithNewNodes(nodeMap)
	cfg, ok := gorumsProvider.StorageConfig()
	assert.True(t, ok)

	responses, readErr := cfg.ReadAll(context.Background(), &kvsprotos.ReadRequest{
		Key: 20,
	})

	// should be as many responses as number of nodes
	assert.Nil(t, readErr)
	assert.Equal(t, len(testNodes)+1, len(responses.GetValue()))
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
	testNodes, wg := oncetree.StartTestNodes(false)

	gorumsProvider := gorumsprovider.New(logger)
	cfg, _ := storageConfig(gorumsProvider)

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

	assert.Nil(t, readErr)
	// should be as many responses as number of nodes
	assert.Equal(t, len(testNodes), len(responses.GetValue()))
	for _, response := range responses.GetValue() {
		assert.Equal(t, int64(100), response)
	}

	for _, node := range testNodes {
		node.Stop("stopped by test")
	}
	wg.Wait()
}

// TestStorageService_WriteLocal tests writing a value to a node, and checking that the LOCAL values has propagated to the required nodes.
func TestStorageService_WriteLocal(t *testing.T) {
	shouldHaveValue := []uint32{1, 2}
	shouldNotHaveValue := []uint32{0, 3, 4, 5, 6, 7, 8, 9}

	testNodes, wg := oncetree.StartTestNodes(false)
	gorumsProvider := gorumsprovider.New(logger)
	cfg, _ := storageConfig(gorumsProvider)

	node0, exits := cfg.Node(0)
	assert.True(t, exits)

	_, writeErr := node0.Write(context.Background(), &kvsprotos.WriteRequest{
		Key:   25,
		Value: 10,
	})
	assert.Nil(t, writeErr)
	time.Sleep(consts.RPCContextTimeout)

	// nodes 1 and 2 should have local value from 0
	for _, gorumsID := range shouldHaveValue {
		node, _ := cfg.Node(gorumsID)
		response, err := node.ReadLocal(
			context.Background(),
			&kvsprotos.ReadLocalRequest{
				Key:    25,
				NodeID: fmt.Sprintf("%d", 0),
			},
		)
		assert.Nil(t, err)
		assert.Equal(t, int64(10), response.GetValue())
	}

	// nodes 0 and 3-9 should not have local value stored from 0
	for _, gorumsID := range shouldNotHaveValue {
		node, _ := cfg.Node(gorumsID)
		response, err := node.ReadLocal(
			context.Background(),
			&kvsprotos.ReadLocalRequest{
				Key:    25,
				NodeID: fmt.Sprintf("%d", 0),
			},
		)
		assert.NotNil(t, err)
		assert.Equal(t, int64(0), response.GetValue())
	}

	for _, node := range testNodes {
		node.Stop("stopped by test")
	}
	wg.Wait()
}

// TestStorageService_RecoverValue crashes a node, then checks if the values are inherited by its parent
// and distributed as a local value to its group
func TestStorageService_RecoverValues(t *testing.T) {
	shouldHaveValue := []uint32{1, 5, 6}
	shouldNotHaveValue := []uint32{0, 3, 4, 7, 8, 9}

	testNodes, wg := oncetree.StartTestNodes(false)
	gorumsProvider := gorumsprovider.New(logger)
	storageCfg, _ := storageConfig(gorumsProvider)
	nodeCfg, _ := nodeConfig(gorumsProvider)

	storageNode2, exists := storageCfg.Node(2)
	assert.True(t, exists)

	ctx, cancel := context.WithTimeout(context.Background(), consts.RPCContextTimeout)
	_, writeErr := storageNode2.Write(ctx, &kvsprotos.WriteRequest{
		Key:   25,
		Value: 10,
	})
	cancel()
	assert.Nil(t, writeErr)
	time.Sleep(consts.RPCContextTimeout)

	// crash node 2
	node2, exists := nodeCfg.Node(2)
	ctx, cancel = context.WithTimeout(context.Background(), consts.RPCContextTimeout)
	_, err := node2.Crash(ctx, &emptypb.Empty{})
	cancel()
	assert.Nil(t, err)
	time.Sleep(consts.FailureDetectionInterval * 3)

	for _, id := range shouldNotHaveValue {
		node, exists := storageCfg.Node(id)
		assert.True(t, exists)
		ctx, cancel = context.WithTimeout(context.Background(), consts.RPCContextTimeout)
		response, err := node.ReadLocal(ctx, &kvsprotos.ReadLocalRequest{
			Key:    25,
			NodeID: "0",
		})
		cancel()
		assert.NotNil(t, err)
		assert.Equal(t, int64(0), response.GetValue())
	}

	for _, id := range shouldHaveValue {
		node, exists := storageCfg.Node(id)
		assert.True(t, exists)
		ctx, cancel = context.WithTimeout(context.Background(), consts.RPCContextTimeout)
		response, err := node.ReadLocal(ctx, &kvsprotos.ReadLocalRequest{
			Key:    25,
			NodeID: "0",
		})
		cancel()
		assert.Nil(t, err)
		assert.Equal(t, int64(10), response.GetValue())
	}

	for _, node := range testNodes {
		node.Stop("stopped by test")
	}
	wg.Wait()
}

func BenchmarkStorageService_Write(t *testing.B) {
	testNodes, wg := oncetree.StartTestNodes(true)
	logger.Info("starting write")
	gorumsProvider := gorumsprovider.New(logger)
	cfg, _ := storageConfig(gorumsProvider)
	nodes := cfg.Nodes()

	t.ResetTimer()
	for i := 0; i < t.N; i++ {
		_, err := nodes[i%len(nodes)].Write(context.Background(), &kvsprotos.WriteRequest{
			Key:   int64(i % 1234),
			Value: int64(i),
		})
		if err != nil {
			panic(err)
		}
	}
	t.StopTimer()
	logger.Info("sleep for rpc context timeout")
	time.Sleep(consts.RPCContextTimeout)
	logger.Info("benchmark done")

	for _, node := range testNodes {
		node.Stop("stopped by test")
	}
	wg.Wait()
}
