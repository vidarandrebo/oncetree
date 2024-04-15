package oncetree

import (
	"context"
	"io"
	"log/slog"
	"net"
	"os"
	"reflect"
	"sync"

	"github.com/vidarandrebo/oncetree/nodemanager/nmenums"

	"github.com/vidarandrebo/oncetree/failuredetector/fdevents"

	"github.com/relab/gorums"
	"github.com/vidarandrebo/oncetree/gorumsprovider"
	fdprotos "github.com/vidarandrebo/oncetree/protos/failuredetector"
	kvsprotos "github.com/vidarandrebo/oncetree/protos/keyvaluestorage"
	"github.com/vidarandrebo/oncetree/protos/node"
	nmprotos "github.com/vidarandrebo/oncetree/protos/nodemanager"

	"github.com/google/uuid"
	"github.com/vidarandrebo/oncetree/consts"
	"github.com/vidarandrebo/oncetree/eventbus"
	"github.com/vidarandrebo/oncetree/failuredetector"
	"github.com/vidarandrebo/oncetree/nodemanager"
	"github.com/vidarandrebo/oncetree/storage"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Node struct {
	id              string
	rpcAddr         string
	gorumsServer    *gorums.Server
	nodeManager     *nodemanager.NodeManager
	logger          *slog.Logger
	timestamp       int64
	failureDetector *failuredetector.FailureDetector
	storageService  *storage.StorageService
	eventbus        *eventbus.EventBus
	stopChan        chan string
}

func NewNode(id string, rpcAddr string, logFile io.Writer) *Node {
	if id == "" {
		id = uuid.New().String()
	}
	logHandlerOpts := slog.HandlerOptions{
		Level: consts.LogLevel,
	}
	logWriter := io.MultiWriter(logFile, os.Stderr)
	logHandler := slog.NewTextHandler(logWriter, &logHandlerOpts)

	logger := slog.New(logHandler).With(slog.Group("node", slog.String("id", id)))
	eventBus := eventbus.New(logger)
	gorumsProvider := gorumsprovider.New(logger)
	nodeManager := nodemanager.New(id, rpcAddr, logger, eventBus, gorumsProvider)
	failureDetector := failuredetector.New(id, logger, nodeManager, eventBus, gorumsProvider)
	storageService := storage.NewStorageService(id, logger, nodeManager, eventBus, gorumsProvider)

	return &Node{
		rpcAddr:         rpcAddr,
		id:              id,
		eventbus:        eventBus,
		nodeManager:     nodeManager,
		logger:          logger.With(slog.Group("node", slog.String("module", "node"))),
		failureDetector: failureDetector,
		storageService:  storageService,
		timestamp:       0,
		stopChan:        make(chan string),
	}
}

func (n *Node) setupEventHandlers() {
	n.eventbus.RegisterHandler(
		reflect.TypeOf(fdevents.NodeFailedEvent{}),
		func(e any) {
			if event, ok := e.(fdevents.NodeFailedEvent); ok {
				n.nodeFailedHandler(event)
			}
		},
	)
}

func (n *Node) nodeFailedHandler(e fdevents.NodeFailedEvent) {
	// n.logger.Printf("node with id %s has failed", e.GroupID)
}

func (n *Node) startGorumsServer(addr string) {
	n.gorumsServer = gorums.NewServer()

	fdprotos.RegisterFailureDetectorServiceServer(n.gorumsServer, n.failureDetector)
	kvsprotos.RegisterKeyValueStorageServer(n.gorumsServer, n.storageService)
	nmprotos.RegisterNodeManagerServiceServer(n.gorumsServer, n.nodeManager)
	node.RegisterNodeServiceServer(n.gorumsServer, n)

	listener, listenErr := net.Listen("tcp", addr)
	if listenErr != nil {
		n.logger.Error(
			"could not create listener",
			slog.String("addr", addr),
			slog.Any("err", listenErr),
		)
		panic("listen failed")
	}
	go func() {
		serveErr := n.gorumsServer.Serve(listener)
		if serveErr != nil {
			n.logger.Error(
				"gorums server could not serve key value server",
				slog.Any("err", serveErr),
			)
			panic("serve failed")
		}
	}()
}

// Run starts the main loop of the node
func (n *Node) Run(knownAddr string) {
	ctx, cancel := context.WithCancel(context.Background())
	wg := sync.WaitGroup{}
	defer cancel()

	n.startGorumsServer(n.rpcAddr)
	n.setupEventHandlers()
	n.nodeManager.SendJoin(knownAddr)

	wg.Add(1)
	go n.failureDetector.Run(ctx, &wg)

	wg.Add(1)
	go n.eventbus.Run(ctx, &wg)

	nodeExitMessage := ""
mainLoop:
	for {
		select {
		case nodeExitMessage = <-n.stopChan:
			n.logger.Info("main loop stopped manually via stop channel")
			cancel()
			break mainLoop
		}
	}
	wg.Wait()
	n.gorumsServer.Stop()
	n.logger.Info("exiting", "msg", nodeExitMessage)
}

func (n *Node) Stop(msg string) {
	n.stopChan <- msg
}

func (n *Node) Crash(ctx gorums.ServerCtx, request *emptypb.Empty) (response *emptypb.Empty, err error) {
	n.logger.Debug("RPC Crash")
	n.Stop("crash RPC")
	return &emptypb.Empty{}, nil
}

// SetNeighboursFromNodeMap assumes a binary tree as slice where a nodes children are at index 2i+1 and 2i+2
// Uses both a slice and a maps to ensure consistent iteration order
func (n *Node) SetNeighboursFromNodeMap(nodeIDs []string, nodes map[string]string) {
	for i, nodeID := range nodeIDs {
		// find n as a child of current node -> current node is n's parent
		if len(nodeIDs) > (2*i+1) && nodeIDs[2*i+1] == n.id {
			n.nodeManager.AddNeighbour(nodeID, nodes[nodeID], nmenums.Parent)
			continue
		}
		if len(nodeIDs) > (2*i+2) && nodeIDs[2*i+2] == n.id {
			n.nodeManager.AddNeighbour(nodeID, nodes[nodeID], nmenums.Parent)
			continue
		}

		// find n -> 2i+1 and 2i+2 are n's children if they exist
		if nodeID == n.id {
			if len(nodeIDs) > (2*i + 1) {
				childId := nodeIDs[2*i+1]
				n.nodeManager.AddNeighbour(childId, nodes[childId], nmenums.Child)
			}
			if len(nodeIDs) > (2*i + 2) {
				childId := nodeIDs[2*i+2]
				n.nodeManager.AddNeighbour(childId, nodes[childId], nmenums.Child)
			}
			continue
		}

	}
	n.logger.Info(
		"tree-info",
		slog.Any("parent", n.nodeManager.Parent()),
	)
	n.logger.Info(
		"tree-info",
		slog.Any("children", n.nodeManager.Children()),
	)
}
