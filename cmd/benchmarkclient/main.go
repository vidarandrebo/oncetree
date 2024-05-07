package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"math/rand"
	"os"
	"time"

	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/vidarandrebo/oncetree/benchmark"

	kvsprotos "github.com/vidarandrebo/oncetree/protos/keyvaluestorage"

	"github.com/vidarandrebo/oncetree/consts"
	nodeprotos "github.com/vidarandrebo/oncetree/protos/node"

	"github.com/vidarandrebo/oncetree/gorumsprovider"
)

func main() {
	// runtime.GOMAXPROCS(1)
	time.Sleep(consts.RPCContextTimeout * 3)
	knownAddr := flag.String("knownAddr", "", "IP address of one of the nodes in the network")
	nodeToCrashAddr := flag.String("nodeToCrashAddr", "", "IP address of one of the node to crash")
	isLeader := flag.Bool("isLeader", false, "Is node that crashes nodes")
	flag.Parse()
	fmt.Println("hello from client")
	fmt.Printf("knownAddr: %s\n", *knownAddr)
	n := 100000

	fmt.Println("isLeader:", *isLeader)

	gorumsProvider := gorumsprovider.New(slog.Default())
	nodeToCrashId := getNodeToCrashID(gorumsProvider, *nodeToCrashAddr)
	fmt.Println("node to crash: ", nodeToCrashId)
	gorumsProvider.Reset()
	benchMarkNodes := mapOnceTreeNodes(gorumsProvider, *knownAddr)

	// all leader will not write to the node that should crash.
	if !*isLeader {
		delete(benchMarkNodes, nodeToCrashId)
	}
	gorumsProvider.ResetWithNewNodes(benchmark.GorumsMap(benchMarkNodes))
	cfg, ok, _ := gorumsProvider.StorageConfig()
	if !ok {
		panic("no storage config")
	}
	WriteStartValues(1000, benchMarkNodes, cfg)
	time.Sleep(consts.RPCContextTimeout)
	fmt.Println("start values written")

	numMsg := int64(0)
	accumulator := 0 * time.Second
	t0 := time.Now()
	timePerRequest := 300 * time.Microsecond
	sleepUnit := 10 * time.Millisecond
	i := 0
	results := make([]benchmark.Result, 0, n+len(benchMarkNodes))
	currentOperationType := benchmark.Read
	for numMsg < 100000 {
		for _, node := range benchMarkNodes {
			// fmt.Println("send request", accumulator, time.Now().Sub(startTime))
			accumulator += timePerRequest
			gorumsNode, ok := cfg.Node(node.GorumsID)
			if !ok {
				panic("no gorums node")
			}
			startTime := time.Now()
			if i%10 == 0 {
				currentOperationType = benchmark.Write
				_, err := gorumsNode.Write(context.Background(), &kvsprotos.WriteRequest{
					Key:   rand.Int63n(1000),
					Value: rand.Int63n(1000000),
				})
				if err != nil {
					fmt.Println("write error")
				}
			} else {
				currentOperationType = benchmark.Read
				_, err := gorumsNode.Read(context.Background(), &kvsprotos.ReadRequest{
					Key: rand.Int63n(1000),
				})
				if err != nil {
					fmt.Println("read error: value not found", err)
				}
			}
			endTime := time.Now()
			results = append(results,
				benchmark.Result{
					Latency:   endTime.Sub(startTime).Microseconds(),
					Timestamp: endTime.UnixMicro(),
					ID:        node.ID,
					Type:      currentOperationType,
				})

			for time.Now().Sub(t0) < accumulator {
				time.Sleep(sleepUnit)
			}
			numMsg++
			if numMsg == 20000 && *isLeader {
				nodeCfg, ok := gorumsProvider.NodeConfig()
				if !ok {
					panic("no gorums node")
				}
				CrashNode(nodeCfg, benchMarkNodes[nodeToCrashId].GorumsID)
				delete(benchMarkNodes, nodeToCrashId)
				gorumsProvider.ResetWithNewNodes(benchmark.GorumsMap(benchMarkNodes))
				cfg, ok, _ = gorumsProvider.StorageConfig()
				break
			}
		}
		i++
	}
	fmt.Println(time.Now().Sub(t0).Milliseconds()-accumulator.Milliseconds(), "ms too slow")
	fmt.Println(results[0])
	id, err := os.Hostname()
	if err != nil {
		panic("failed to get hostname")
	}
	benchmark.WriteResultsToDisk(results, id)
}

func CrashNode(cfg *nodeprotos.Configuration, gorumsID uint32) {
	node, ok := cfg.Node(gorumsID)
	if !ok {
		panic("no gorums node")
	}
	_, err := node.Crash(context.Background(), &emptypb.Empty{})
	if err != nil {
		panic("failed to crash node")
	}
	fmt.Println("crashed node")
}

func WriteStartValues(n int64, benchMarkNodes map[string]benchmark.Node, cfg *kvsprotos.Configuration) {
	for i := int64(0); i < n; i++ {
		for _, node := range benchMarkNodes {
			gorumsNode, ok := cfg.Node(node.GorumsID)
			if !ok {
				panic("no gorums node")
			}
			_, err := gorumsNode.Write(context.Background(), &kvsprotos.WriteRequest{
				Key:   i,
				Value: rand.Int63n(1000000),
			})
			if err != nil {
				panic(fmt.Sprintf("write error: %v", err))
			}
			time.Sleep(1 * time.Millisecond)
		}
	}
}

func mapOnceTreeNodes(provider *gorumsprovider.GorumsProvider, knownAddr string) map[string]benchmark.Node {
	provider.SetNodes(map[string]uint32{
		knownAddr: 0,
	})
	id, err := os.Hostname()
	if err != nil {
		panic("failed to get hostname: " + err.Error())
	}
	cfg, ok := provider.NodeConfig()
	if !ok {
		panic("failed to get node config")
	}
	node, ok := cfg.Node(0)
	if !ok {
		panic("failed to get node")
	}
	ctx, cancel := context.WithTimeout(context.Background(), consts.RPCContextTimeout)
	defer cancel()
	response, err := node.Nodes(ctx, &nodeprotos.NodesRequest{Origin: id})
	if err != nil {
		panic("failed to get nodes")
	}
	result := make(map[string]benchmark.Node)
	nodeMap := response.GetNodeMap()
	gorumsID := uint32(0)
	for id, address := range nodeMap {
		newNode := benchmark.Node{
			ID:       id,
			Address:  address,
			GorumsID: gorumsID,
		}
		result[id] = newNode
		gorumsID = gorumsID + 1
	}
	return result
}

func getNodeToCrashID(provider *gorumsprovider.GorumsProvider, nodeToCrashAddr string) string {
	provider.SetNodes(map[string]uint32{
		nodeToCrashAddr: 0,
	})
	cfg, ok := provider.NodeConfig()
	if !ok {
		panic("failed to get node config")
	}
	node, ok := cfg.Node(0)
	if !ok {
		panic("failed to get node")
	}
	ctx, cancel := context.WithTimeout(context.Background(), consts.RPCContextTimeout)
	defer cancel()
	response, err := node.NodeID(ctx, &emptypb.Empty{})
	if err != nil {
		panic("failed to get id of node to crash")
	}
	return response.GetID()
}
