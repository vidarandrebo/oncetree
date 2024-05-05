package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"math/rand"
	"os"
	"runtime"
	"time"

	kvsprotos "github.com/vidarandrebo/oncetree/protos/keyvaluestorage"

	"github.com/vidarandrebo/oncetree/consts"
	nodeprotos "github.com/vidarandrebo/oncetree/protos/node"

	"github.com/vidarandrebo/oncetree/gorumsprovider"
)

func main() {
	runtime.GOMAXPROCS(1)
	time.Sleep(consts.RPCContextTimeout * 3)
	knownAddr := flag.String("knownAddr", "", "IP address of one of the nodes in the network")
	flag.Parse()
	fmt.Println("hello from client")
	fmt.Printf("knownAddr: %s\n", *knownAddr)

	gorumsProvider := gorumsprovider.New(slog.Default())
	benchMarkNodes := mapOnceTreeNodes(gorumsProvider, *knownAddr)
	gorumsProvider.ResetWithNewNodes(GorumsMap(benchMarkNodes))
	fmt.Println(benchMarkNodes)
	fmt.Println(len(benchMarkNodes))
	cfg, ok, _ := gorumsProvider.StorageConfig()
	if !ok {
		panic("no storage config")
	}
	numMsg := int64(0)
	for {
		for _, node := range benchMarkNodes {
			gorumsNode, ok := cfg.Node(node.GorumsID)
			if !ok {
				panic("no gorums node")
			}
			if numMsg%10 == 0 || numMsg < 5000 {
				_, err := gorumsNode.Write(context.Background(), &kvsprotos.WriteRequest{
					Key:   rand.Int63n(1000),
					Value: rand.Int63n(1000),
				})
				// cancel()
				if err != nil {
					panic(fmt.Sprintf("write error: %v", err))
				}
			} else {
				_, err := gorumsNode.Read(context.Background(), &kvsprotos.ReadRequest{
					Key: rand.Int63n(1000),
				})
				if err != nil {
					fmt.Println("read error: value not found")
				}
			}
			numMsg++
			if numMsg%10000 == 0 {
				fmt.Println("sent", numMsg, "messages")
			}
		}
		// time.Sleep(10 * time.Millisecond)
	}
}

type BenchMarkNode struct {
	ID       string
	Address  string
	GorumsID uint32
}

func (bn BenchMarkNode) String() string {
	return fmt.Sprintf("ID: %s, Address: %s, GorumsID %d", bn.ID, bn.Address, bn.GorumsID)
}

func GorumsMap(nodes map[string]BenchMarkNode) map[string]uint32 {
	gorumsMap := make(map[string]uint32)
	for _, node := range nodes {
		gorumsMap[node.Address] = node.GorumsID
	}
	return gorumsMap
}

func mapOnceTreeNodes(provider *gorumsprovider.GorumsProvider, knownAddr string) map[string]BenchMarkNode {
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
	result := make(map[string]BenchMarkNode)
	nodeMap := response.GetNodeMap()
	gorumsID := uint32(0)
	for id, address := range nodeMap {
		newNode := BenchMarkNode{
			ID:       id,
			Address:  address,
			GorumsID: gorumsID,
		}
		result[id] = newNode
		gorumsID = gorumsID + 1
	}
	return result
}
