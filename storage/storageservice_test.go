//go:build !race

package storage_test

import (
	"context"
	"fmt"
	"log/slog"
	"maps"
	"math"
	"math/rand"
	"os"
	"slices"
	"strconv"
	"testing"
	"time"

	"github.com/vidarandrebo/oncetree/concurrent/hashset"

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

func generateRandomWrites(n int) map[string]map[int64]int64 {
	// rand.Seed(0)
	values := make(map[string]map[int64]int64)
	for _, id := range nodeIDs {
		values[id] = make(map[int64]int64)
		for i := 0; i < n; i++ {
			values[id][rand.Int63n(int64(n))] = rand.Int63n(math.MaxInt32)
		}
	}
	return values
}

func keysFromValueMap(valueMap map[string]map[int64]int64) hashset.HashSet[int64] {
	hashSet := hashset.NewHashSet[int64]()
	for _, values := range valueMap {
		for key := range values {
			hashSet.Add(key)
		}
	}
	return hashSet
}

func keysFromValueMapForNodes(nodes []string, valueMap map[string]map[int64]int64) hashset.HashSet[int64] {
	hashSet := hashset.NewHashSet[int64]()
	for nodeID, values := range valueMap {
		if slices.Contains(nodes, nodeID) {
			for key := range values {
				hashSet.Add(key)
			}
		}
	}
	return hashSet
}

func aggValueForKey(key int64, allValue map[string]map[int64]int64) int64 {
	sum := int64(0)
	for _, values := range allValue {
		value, ok := values[key]
		if ok {
			sum += value
		}
	}
	return sum
}

func nodesValueForKey(key int64, nodes []string, allValues map[string]map[int64]int64) int64 {
	sum := int64(0)
	for id, values := range allValues {
		value, ok := values[key]
		if ok && slices.Contains(nodes, id) {
			sum += value
		}
	}
	return sum
}

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
	cfg, _, _ := provider.StorageConfig()
	return cfg, nodeMap
}

func writeRandomValuesToNodes(cfg *kvsprotos.Configuration, n int) map[string]map[int64]int64 {
	valuesToWrite := generateRandomWrites(n)
	for nodeID, values := range valuesToWrite {
		gorumsID, _ := strconv.Atoi(nodeID)
		node, ok := cfg.Node(uint32(gorumsID))
		if !ok {
			panic("failed to find node in config")
		}
		for key, value := range values {
			ctx, cancel := context.WithTimeout(context.Background(), consts.RPCContextTimeout)
			_, err := node.Write(ctx, &kvsprotos.WriteRequest{
				Key:   key,
				Value: value,
			})
			cancel()
			if err != nil {
				panic(fmt.Sprintf("failed to write value to node %s", nodeID))
			}
		}
	}
	return valuesToWrite
}

// TestStorageService_shareAll tests writing the values to all nodes, and checking that the values have propagated to all nodes.
func TestStorageService_shareAll(t *testing.T) {
	testNodes, wg := oncetree.StartTestNodes(true)

	gorumsProvider := gorumsprovider.New(logger)
	cfg, nodeMap := storageConfig(gorumsProvider)

	writtenValues := writeRandomValuesToNodes(cfg, 100)
	time.Sleep(consts.TestWaitAfterWrite)
	keys := keysFromValueMap(writtenValues)
	newNode, newWg := oncetree.StartTestNode(true)

	// add the new node to config
	nodeMap[":9090"] = 10
	gorumsProvider.ResetWithNewNodes(nodeMap)
	cfg, ok, _ := gorumsProvider.StorageConfig()
	assert.True(t, ok)
	time.Sleep(consts.TestWaitAfterWrite)

	for _, key := range keys.Values() {
		responses, readErr := cfg.ReadAll(context.Background(), &kvsprotos.ReadRequest{
			Key: key,
		})
		// should be as many responses as number of nodes
		assert.Nil(t, readErr)
		assert.Equal(t, len(testNodes)+1, len(responses.GetValue()))
		// checks that all nodes including the new has the value 100 for key 20
		shouldBe := aggValueForKey(key, writtenValues)
		for _, actual := range responses.GetValue() {
			assert.Equal(t, shouldBe, actual)
		}
	}

	newNode.Stop("stopped by test")
	for _, node := range testNodes {
		node.Stop("stopped by test")
	}
	wg.Wait()
	newWg.Wait()
}

// TestStorageService_Write tests writing to all nodes and check that the values have propagated to all nodes
func TestStorageService_Write(t *testing.T) {
	testNodes, wg := oncetree.StartTestNodes(false)

	gorumsProvider := gorumsprovider.New(logger)
	cfg, _ := storageConfig(gorumsProvider)

	writtenValues := writeRandomValuesToNodes(cfg, 100)
	time.Sleep(consts.TestWaitAfterWrite)
	keys := keysFromValueMap(writtenValues)

	for _, key := range keys.Values() {
		ctx := context.Background()
		responses, err := cfg.ReadAll(ctx, &kvsprotos.ReadRequest{
			Key: key,
		})
		assert.Nil(t, err)
		assert.Equal(t, len(testNodes), len(responses.GetValue()))
		shouldBe := aggValueForKey(key, writtenValues)
		for _, actual := range responses.GetValue() {
			assert.Equal(t, shouldBe, actual)
		}
	}

	for _, node := range testNodes {
		node.Stop("stopped by test")
	}
	wg.Wait()
}

// TestStorageService_WriteLocal tests writing random values to all nodes, and checking that the LOCAL values has propagated to the required nodes.
/*
				0
               / \
              1   2
            / |   | \
           3  4   5  6
         / |   \
		7  8    9
*/
func TestStorageService_WriteLocal(t *testing.T) {
	shouldHaveValue := []uint32{0, 5, 6}
	shouldNotHaveValue := []uint32{1, 2, 3, 4, 7, 8, 9}

	testNodes, wg := oncetree.StartTestNodes(false)
	gorumsProvider := gorumsprovider.New(logger)
	storageCfg, _ := storageConfig(gorumsProvider)

	writtenValues := writeRandomValuesToNodes(storageCfg, 100)
	time.Sleep(consts.TestWaitAfterWrite)
	nodes := []string{"2"}
	keys := keysFromValueMapForNodes(nodes, writtenValues)

	for _, key := range keys.Values() {
		for _, id := range shouldNotHaveValue {
			node, exists := storageCfg.Node(id)
			assert.True(t, exists)
			ctx, cancel := context.WithTimeout(context.Background(), consts.RPCContextTimeout)
			response, err := node.ReadLocal(ctx, &kvsprotos.ReadLocalRequest{
				Key:    key,
				NodeID: "2",
			})
			cancel()
			assert.NotNil(t, err)
			assert.Equal(t, int64(0), response.GetValue())
		}

		for _, id := range shouldHaveValue {
			node, exists := storageCfg.Node(id)
			assert.True(t, exists)
			ctx, cancel := context.WithTimeout(context.Background(), consts.RPCContextTimeout)
			response, err := node.ReadLocal(ctx, &kvsprotos.ReadLocalRequest{
				Key:    key,
				NodeID: "2",
			})
			cancel()
			assert.Nil(t, err)
			assert.Equal(t, nodesValueForKey(key, nodes, writtenValues), response.GetValue())
		}
	}

	for _, node := range testNodes {
		node.Stop("stopped by test")
	}
	wg.Wait()
}

// TestStorageService_RecoverValue crashes a node, then checks if the values are inherited by its parent
// and distributed as a local value to its group
/*
				0
               / \
              1   2
            / |   | \
           3  4   5  6
         / |   \
		7  8    9
*/
// Will be converted to:
/*
				 0
               / | \
              1  5  6
            / |
           3  4
         / |   \
		7  8    9
*/
func TestStorageService_RecoverValues(t *testing.T) {
	shouldHaveValue := []uint32{1, 5, 6}
	shouldNotHaveValue := []uint32{0, 3, 4, 7, 8, 9}

	testNodes, wg := oncetree.StartTestNodes(false)
	gorumsProvider := gorumsprovider.New(logger)
	storageCfg, _ := storageConfig(gorumsProvider)
	nodeCfg, _ := nodeConfig(gorumsProvider)

	valuesWritten := writeRandomValuesToNodes(storageCfg, 100)
	joinedNodes := []string{"0", "2"}
	keys := keysFromValueMapForNodes(joinedNodes, valuesWritten)

	time.Sleep(consts.TestWaitAfterWrite)

	// crash node 2
	node2, exists := nodeCfg.Node(2)
	assert.True(t, exists)
	ctx := context.Background()
	_, err := node2.Crash(ctx, &emptypb.Empty{})
	assert.Nil(t, err)
	time.Sleep(consts.HeartbeatSendInterval * consts.FailureDetectorStrikes * 2)

	for _, key := range keys.Values() {
		for _, id := range shouldNotHaveValue {
			node, exists := storageCfg.Node(id)
			assert.True(t, exists)
			ctx = context.Background()
			response, err := node.ReadLocal(ctx, &kvsprotos.ReadLocalRequest{
				Key:    key,
				NodeID: "0",
			})
			assert.NotNil(t, err)
			assert.Equal(t, int64(0), response.GetValue())
		}

		for _, id := range shouldHaveValue {
			node, exists := storageCfg.Node(id)
			assert.True(t, exists)
			ctx = context.Background()
			response, err := node.ReadLocal(ctx, &kvsprotos.ReadLocalRequest{
				Key:    key,
				NodeID: "0",
			})
			assert.Nil(t, err)
			assert.Equal(t, nodesValueForKey(key, joinedNodes, valuesWritten), response.GetValue())
		}
	}

	for _, node := range testNodes {
		node.Stop("stopped by test")
	}
	wg.Wait()
}

func BenchmarkStorageService_Write(t *testing.B) {
	// runtime.GOMAXPROCS(2)
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
