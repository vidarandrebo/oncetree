package oncetree

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"log"
	"net"
	"os"
	"reflect"
	"sync"
	"time"

	"github.com/vidarandrebo/oncetree/eventbus"
	"github.com/vidarandrebo/oncetree/failuredetector"
	"github.com/vidarandrebo/oncetree/nodemanager"
	"github.com/vidarandrebo/oncetree/storage"

	"github.com/relab/gorums"
	fdprotos "github.com/vidarandrebo/oncetree/protos/failuredetector"
	kvsprotos "github.com/vidarandrebo/oncetree/protos/keyvaluestorage"
	"github.com/vidarandrebo/oncetree/protos/node"
	nmprotos "github.com/vidarandrebo/oncetree/protos/nodemanager"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Node struct {
	id                     string
	rpcAddr                string
	gorumsServer           *gorums.Server
	gorumsManagers         *GorumsManagers
	nodeManager            *nodemanager.NodeManager
	logger                 *log.Logger
	timestamp              int64
	failureDetector        *failuredetector.FailureDetector
	keyValueStorageService *storage.KeyValueStorageService
	eventbus               *eventbus.EventBus
	stopChan               chan string
}

func NewNode(id string, rpcAddr string) *Node {
	if id == "" {
		id = uuid.New().String()
	}
	logger := log.New(os.Stderr, fmt.Sprintf("NodeID: %s ", id), log.Ltime|log.Lmsgprefix)
	nodeManager := nodemanager.New(id, 2)
	eventBus := eventbus.New(logger)
	gorumsManagers := CreateGorumsManagers()
	return &Node{
		rpcAddr: rpcAddr,
		id:      id,
		failureDetector: failuredetector.New(
			id,
			logger,
			nodeManager,
			gorumsManagers.fdManager,
			eventBus,
		),
		keyValueStorageService: storage.NewKeyValueStorageService(
			id,
			logger,
			nodeManager,
			gorumsManagers.kvsManager,
		),
		eventbus:       eventBus,
		nodeManager:    nodeManager,
		gorumsManagers: gorumsManagers,
		logger:         logger,
		timestamp:      0,
		stopChan:       make(chan string),
	}
}

// setupGorumsConfigs sets up the initial gorums configurations for the various services
func (n *Node) setupGorumsConfigs() {
	tries := 0
startTryCreateConfigs:
	time.Sleep(time.Duration(tries))
	tries++
	if tries > 5 {
		log.Panicln("could not create gorums configs")
	}
	fdSetupErr := n.failureDetector.SetNodesFromManager()
	kvssSetupErr := n.keyValueStorageService.SetNodesFromManager()
	if fdSetupErr != nil {
		n.logger.Println(fdSetupErr)
		goto startTryCreateConfigs
	}
	if kvssSetupErr != nil {
		n.logger.Println(kvssSetupErr)
		goto startTryCreateConfigs
	}
}

func (n *Node) setupEventHandlers() {
	n.eventbus.RegisterHandler(reflect.TypeOf(failuredetector.NodeFailedEvent{}), func(e any) {
		if event, ok := e.(failuredetector.NodeFailedEvent); ok {
			n.nodeFailedHandler(event)
		}
	})
}

func (n *Node) nodeFailedHandler(e failuredetector.NodeFailedEvent) {
	n.logger.Printf("node with id %s has failed", e.NodeID)
}

func (n *Node) startGorumsServer(addr string) {
	n.gorumsServer = gorums.NewServer()

	fdprotos.RegisterFailureDetectorServiceServer(n.gorumsServer, n.failureDetector)
	kvsprotos.RegisterKeyValueStorageServer(n.gorumsServer, n.keyValueStorageService)
	nmprotos.RegisterNodeManagerServiceServer(n.gorumsServer, n.nodeManager)
	node.RegisterNodeServiceServer(n.gorumsServer, n)

	listener, listenErr := net.Listen("tcp", addr)
	if listenErr != nil {
		n.logger.Panicf("could not listen to address %v", addr)
	}
	go func() {
		serveErr := n.gorumsServer.Serve(listener)
		if serveErr != nil {
			n.logger.Panicf("gorums server could not serve key value server")
		}
	}()
}

// Run starts the main loop of the node
func (n *Node) Run() {
	ctx, cancel := context.WithCancel(context.Background())
	wg := sync.WaitGroup{}
	defer cancel()
	n.startGorumsServer(n.rpcAddr)
	n.setupGorumsConfigs()
	n.setupEventHandlers()
	wg.Add(1)
	go n.failureDetector.Run(ctx, &wg)

	nodeExitMessage := ""
	n.logger.Println("hello there")
mainLoop:
	for {
		select {
		case nodeExitMessage = <-n.stopChan:
			n.logger.Println("Main loop stopped manually via stop channel")
			cancel()
			break mainLoop
		}
	}
	wg.Wait()
	n.gorumsServer.Stop()
	n.logger.Printf("Exiting after message \"%s\" on stop channel", nodeExitMessage)
}

func (n *Node) Stop(msg string) {
	n.stopChan <- msg
}

func (n *Node) Crash(ctx gorums.ServerCtx, request *emptypb.Empty) (response *emptypb.Empty, err error) {
	n.Stop("crash RPC")
	return &emptypb.Empty{}, nil
}

// SetNeighboursFromNodeMap assumes a binary tree as slice where a nodes children are at index 2i+1 and 2i+2
// Uses both a slice and a maps to ensure consistent iteration order
func (n *Node) SetNeighboursFromNodeMap(nodeIDs []string, nodes map[string]string) {
	for i, nodeID := range nodeIDs {
		// find n as a child of current node -> current node is n's parent
		if len(nodeIDs) > (2*i+1) && nodeIDs[2*i+1] == n.id {
			n.nodeManager.AddNeighbour(nodeID, nodes[nodeID], nodemanager.Parent)
			continue
		}
		if len(nodeIDs) > (2*i+2) && nodeIDs[2*i+2] == n.id {
			n.nodeManager.AddNeighbour(nodeID, nodes[nodeID], nodemanager.Parent)
			continue
		}

		// find n -> 2i+1 and 2i+2 are n's children if they exist
		if nodeID == n.id {
			if len(nodeIDs) > (2*i + 1) {
				childId := nodeIDs[2*i+1]
				n.nodeManager.AddNeighbour(childId, nodes[childId], nodemanager.Child)
			}
			if len(nodeIDs) > (2*i + 2) {
				childId := nodeIDs[2*i+2]
				n.nodeManager.AddNeighbour(childId, nodes[childId], nodemanager.Child)
			}
			continue
		}

	}
	n.logger.Printf("parent: %v", n.nodeManager.GetParent())
	n.logger.Printf("children: %v", n.nodeManager.GetChildren())
}
