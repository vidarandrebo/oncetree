package storage_test

import (
	"context"
	"fmt"
	"log/slog"
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
	logger = slog.Default().With(slog.Group("node", slog.String("module", "storage_test")))
)

// TestStorageService_Write tests writing the same value to all nodes, and checking that the values has propagated to all nodes.
func TestStorageService_shareAll(t *testing.T) {
	testNodes, wg := oncetree.StartTestNodes(true)
	time.Sleep(consts.GorumsDialTimeout)
	mgr, cfg := createKeyValueStorageConfig()

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
	testNodes, wg := oncetree.StartTestNodes(false)
	time.Sleep(consts.GorumsDialTimeout)
	_, cfg := createKeyValueStorageConfig()

	for i, node := range cfg.Nodes() {
		_, writeErr := node.Write(context.Background(), &kvsprotos.WriteRequest{
			Key:     20,
			Value:   10,
			WriteID: int64(i + 1),
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

// TestStorageService_WriteLocal tests writing a value to a node, and checking that the LOCAL values has propagated to the required nodes.
func TestStorageService_WriteLocal(t *testing.T) {
	shouldHaveValue := []uint32{1, 2}
	shouldNotHaveValue := []uint32{0, 3, 4, 5, 6, 7, 8, 9}
	testNodes, wg := oncetree.StartTestNodes(false)
	time.Sleep(consts.GorumsDialTimeout)
	_, cfg := createKeyValueStorageConfig()
	node_0, exits := cfg.Node(0)
	assert.True(t, exits)

	_, writeErr := node_0.Write(context.Background(), &kvsprotos.WriteRequest{
		Key:     25,
		Value:   10,
		WriteID: int64(1),
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

// createKeyValueStorageConfig creates a new manager and returns an initialized configuration to use with the StorageService
func createKeyValueStorageConfig() (*kvsprotos.Manager, *kvsprotos.Configuration) {
	manager := kvsprotos.NewManager(
		gorums.WithDialTimeout(consts.GorumsDialTimeout),
		gorums.WithGrpcDialOptions(
			grpc.WithBlock(),
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		),
	)
	cfg, err := manager.NewConfiguration(&sqspec.QSpec{NumNodes: len(nodeAddrs)}, gorums.WithNodeMap(gorumsNodeMap))
	if err != nil {
		panic("failed to create cfg")
	}
	return manager, cfg
}

func BenchmarkStorageService_Write(t *testing.B) {
	testNodes, wg := oncetree.StartTestNodes(true)
	logger.Info("sleep for dial timeout before starting to write")
	time.Sleep(consts.RPCContextTimeout * 2)
	logger.Info("starting write")
	_, cfg := createKeyValueStorageConfig()
	nodes := cfg.Nodes()

	for i := 0; i < t.N; i++ {
		_, err := nodes[i%len(nodes)].Write(context.Background(), &kvsprotos.WriteRequest{
			Key:     int64(i % 1234),
			Value:   int64(i),
			WriteID: int64(i),
		})
		if err != nil {
			panic(err)
		}
		// log.Printf("[Bench] %d", i)
	}
	logger.Info("sleep for rpc context timeout")
	time.Sleep(consts.RPCContextTimeout)
	logger.Info("benchmark done")

	for _, node := range testNodes {
		node.Stop("stopped by test")
	}
	wg.Wait()
}
