package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log/slog"
	"math/rand"
	"os"
	"path/filepath"
	"time"

	"github.com/google/uuid"

	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/vidarandrebo/oncetree/benchmark"

	kvsprotos "github.com/vidarandrebo/oncetree/protos/keyvaluestorage"

	"github.com/vidarandrebo/oncetree/consts"
	nodeprotos "github.com/vidarandrebo/oncetree/protos/node"

	"github.com/vidarandrebo/oncetree/gorumsprovider"
)

func main() {
	id, err := os.Hostname()
	if err != nil {
		panic("could not get hostname")
	}
	id += "_"
	id += uuid.New().String()
	fileName := filepath.Join(consts.LogFolder, fmt.Sprintf("client_%s.log", id))
	file, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o666)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	logHandlerOpts := slog.HandlerOptions{
		Level: consts.LogLevel,
	}
	logWriter := io.MultiWriter(file, os.Stderr)
	logHandler := slog.NewTextHandler(logWriter, &logHandlerOpts)

	logger := slog.New(logHandler).With(slog.Group("client", slog.String("id", id)))

	// runtime.GOMAXPROCS(1)
	time.Sleep(consts.RPCContextTimeout * 3)
	knownAddr := flag.String("known-address", "", "IP address of one of the nodes in the network")
	writer := flag.Bool("writer", false, "client is writer")
	reader := flag.Bool("reader", false, "client is reader")
	nodeToCrashAddr := flag.String("node-to-crash-address", "", "IP address of one of the node to crash")
	leader := flag.Bool("leader", false, "Is node that crashes nodes")
	flag.Parse()
	logger.Info("hello from client")
	logger.Info("known-address", slog.String("value", *knownAddr))
	n := 100000

	fmt.Println("is-leader:", *leader)

	gorumsProvider := gorumsprovider.New(slog.Default())
	nodeToCrashId := getNodeToCrashID(gorumsProvider, *nodeToCrashAddr)
	fmt.Println("node to crash: ", nodeToCrashId)
	gorumsProvider.Reset()
	benchMarkNodes := mapOnceTreeNodes(gorumsProvider, *knownAddr)

	gorumsProvider.ResetWithNewNodes(benchmark.GorumsMap(benchMarkNodes))
	storageCfg, ok, _ := gorumsProvider.StorageConfig()
	if !ok {
		panic("no storage config")
	}
	WriteStartValues(consts.BenchmarkNumKeys, benchMarkNodes, storageCfg, logger)
	time.Sleep(2 * consts.RPCContextTimeout)
	logger.Info("start values written")

	numMsg := int64(0)
	accumulator := 0 * time.Second
	t0 := time.Now()
	timePerRequest := 1000 * time.Microsecond
	sleepUnit := 10 * time.Millisecond
	i := 0
	hasCrashed := false
	results := make([]benchmark.Result, 0, n+len(benchMarkNodes))
	operationType := benchmark.Read
	if *writer {
		logger.Info("client is writer")
		operationType = benchmark.Write
	} else if *reader {
		logger.Info("client is reader")
		operationType = benchmark.Read
	}
	clients := make([]benchmark.IClient, 0)
	for _, node := range benchMarkNodes {
		cfg, err := gorumsProvider.CustomStorageConfig(map[string]uint32{
			node.Address: node.GorumsID,
		})
		if err != nil {
			panic(err)
		}
		if *writer {
			clients = append(clients, benchmark.NewWriter(cfg, 1000, node.ID))
		} else {
			clients = append(clients, benchmark.NewReader(cfg, 1000, node.ID))
		}
	}
	ctx, cancel := context.WithCancel(context.Background())
	for _, client := range clients {
		client.Run(ctx)
	}
	numSkips := 0
	for time.Now().Sub(t0) < consts.BenchmarkTime {
		for _, client := range clients {
			started := client.Send()
			if !started {
				numSkips++
				if numSkips > len(benchMarkNodes) {
					time.Sleep(sleepUnit)
				}
				// logger.Info("skipping")
				continue
			}
			numSkips = 0
			accumulator += timePerRequest

			numMsg++
			if *leader && !hasCrashed && time.Now().Sub(t0) > 20*time.Second {
				nodeCfg, ok := gorumsProvider.NodeConfig()
				if !ok {
					panic("no gorums node")
				}
				CrashNode(nodeCfg, benchMarkNodes[nodeToCrashId].GorumsID)
				hasCrashed = true
			}
			for time.Now().Sub(t0) < accumulator {
				time.Sleep(sleepUnit)
			}
		}
		i++
	}
	logger.Info("time to slow",
		slog.Int64("ms", time.Now().Sub(t0).Milliseconds()-accumulator.Milliseconds()))
	time.Sleep(time.Second)
	cancel()
	for _, client := range clients {
		clientResult := client.GetResults()
		logger.Info("join client", slog.Int("num msg", len(clientResult)))
		results = append(results, clientResult...)
	}
	fmt.Println(results[0])
	benchmark.WriteResultsToDisk(results, id, operationType)
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

func WriteStartValues(n int64, benchMarkNodes map[string]benchmark.Node, cfg *kvsprotos.Configuration, logger *slog.Logger) {
	accumulator := 0 * time.Second
	// 25 % writers -> 4 times slower when readers are writing
	timePerRequest := 4 * 1000 * time.Microsecond
	t0 := time.Now()
	sleepUnit := 10 * time.Millisecond
	for i := int64(0); i < n; i++ {
		for _, node := range benchMarkNodes {
			accumulator += timePerRequest
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
			for time.Now().Sub(t0) < accumulator {
				time.Sleep(sleepUnit)
			}
		}
	}
	logger.Info("time to slow start values",
		slog.Int64("ms", time.Now().Sub(t0).Milliseconds()-accumulator.Milliseconds()))
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
	ctx := context.Background()
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
	ctx := context.Background()
	response, err := node.NodeID(ctx, &emptypb.Empty{})
	if err != nil {
		panic("failed to get id of node to crash")
	}
	return response.GetID()
}
