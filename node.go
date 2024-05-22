package oncetree

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"net"
	"os"
	"reflect"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/relab/gorums"
	"github.com/vidarandrebo/oncetree/consts"
	"github.com/vidarandrebo/oncetree/eventbus"
	"github.com/vidarandrebo/oncetree/failuredetector"
	"github.com/vidarandrebo/oncetree/failuredetector/fdevents"
	"github.com/vidarandrebo/oncetree/gorumsprovider"
	"github.com/vidarandrebo/oncetree/nodemanager"
	fdprotos "github.com/vidarandrebo/oncetree/protos/failuredetector"
	kvsprotos "github.com/vidarandrebo/oncetree/protos/keyvaluestorage"
	nodeprotos "github.com/vidarandrebo/oncetree/protos/node"
	nmprotos "github.com/vidarandrebo/oncetree/protos/nodemanager"
	"github.com/vidarandrebo/oncetree/storage"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Node struct {
	id              string
	rpcAddr         string
	rpcPort         string
	gorumsServer    *gorums.Server
	nodeManager     *nodemanager.NodeManager
	logger          *slog.Logger
	timestamp       int64
	failureDetector *failuredetector.FailureDetector
	storageService  *storage.StorageService
	eventbus        *eventbus.EventBus
	stopChan        chan string
}

func (n *Node) NodeID(ctx gorums.ServerCtx, request *emptypb.Empty) (response *nodeprotos.IDResponse, err error) {
	return &nodeprotos.IDResponse{ID: n.id}, nil
}

func NewNode(id string, rpcAddr string, rpcPort string, logFile io.Writer) *Node {
	if id == "" {
		id = uuid.New().String()
	}
	logHandlerOpts := slog.HandlerOptions{
		Level: consts.LogLevel,
	}
	logWriter := io.MultiWriter(logFile, os.Stderr)
	//	if logFile == io.Discard {
	//		logWriter = io.Discard
	//	}
	logHandler := slog.NewTextHandler(logWriter, &logHandlerOpts)

	address := fmt.Sprintf("%s:%s", rpcAddr, rpcPort)

	logger := slog.New(logHandler).With(slog.Group("node", slog.String("id", id)))
	eventBus := eventbus.New(logger)
	gorumsProvider := gorumsprovider.New(logger)
	nodeManager := nodemanager.New(id, address, logger, eventBus, gorumsProvider)
	failureDetector := failuredetector.New(id, logger, nodeManager, eventBus, gorumsProvider)
	storageService := storage.NewStorageService(id, logger, nodeManager, eventBus, gorumsProvider)
	ips, err := net.LookupIP(strings.Split(rpcAddr, ":")[0])
	if err == nil {
		for _, ip := range ips {
			logger.Info(ip.String())
		}
	} else {
		logger.Error("failed to lookup ips")
	}

	return &Node{
		id:              id,
		rpcAddr:         rpcAddr,
		rpcPort:         rpcPort,
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
	n.logger.Info("serving at",
		slog.String("addr", addr))
	n.gorumsServer = gorums.NewServer()

	fdprotos.RegisterFailureDetectorServiceServer(n.gorumsServer, n.failureDetector)
	kvsprotos.RegisterKeyValueStorageServer(n.gorumsServer, n.storageService)
	nmprotos.RegisterNodeManagerServiceServer(n.gorumsServer, n.nodeManager)
	nodeprotos.RegisterNodeServiceServer(n.gorumsServer, n)

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
func (n *Node) Run(knownAddr string, readyCallBack func()) {
	ctx, cancel := context.WithCancel(context.Background())
	wg := sync.WaitGroup{}
	defer cancel()

	// serve on all addresses
	serveAddr := fmt.Sprintf("%s:%s", "", n.rpcPort)

	n.startGorumsServer(serveAddr)
	n.setupEventHandlers()
	n.nodeManager.SendJoin(knownAddr)

	wg.Add(1)
	go n.failureDetector.Run(ctx, &wg)

	wg.Add(1)
	go n.eventbus.Run(ctx, &wg)

	nodeExitMessage := ""
	if readyCallBack != nil {
		readyCallBack()
	}
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
	select {
	case n.stopChan <- msg:
		n.logger.Info("pushed to stop chan")
	default:
		n.logger.Info("no log chan, exiting")

	}
}

func (n *Node) Crash(ctx gorums.ServerCtx, request *emptypb.Empty) (response *emptypb.Empty, err error) {
	n.logger.Info("RPC Crash")
	go func() {
		time.Sleep(time.Second)
		n.Stop("crash RPC")
	}()
	return &emptypb.Empty{}, nil
}

func (n *Node) Nodes(ctx gorums.ServerCtx, request *nodeprotos.NodesRequest) (*nodeprotos.NodesResponse, error) {
	nodeMap := make(map[string]string)

	nodeMap[n.id] = fmt.Sprintf("%s:%s", n.rpcAddr, n.rpcPort)
	neighbours := n.nodeManager.Neighbours()
	cfg, ok := n.nodeManager.NodeConfig()
	if !ok {
		return nil, fmt.Errorf("node config not found")
	}
	for _, neighbour := range neighbours {
		if neighbour.Key == request.Origin {
			continue
		}
		gorumsNode, ok := cfg.Node(neighbour.Value.GorumsID)
		if !ok {
			return nil, fmt.Errorf("neighbour %s not found in node config", neighbour.Value)
		}
		response, err := gorumsNode.Nodes(ctx, &nodeprotos.NodesRequest{Origin: n.id})
		if err != nil {
			return nil, err
		}
		for id, address := range response.GetNodeMap() {
			nodeMap[id] = address
		}
	}
	response := nodeprotos.NodesResponse{NodeMap: nodeMap}
	return &response, nil
}
