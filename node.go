package oncetree

import (
	"context"
	"io"
	"log/slog"
	"net"
	"os"
	"reflect"
	"sync"

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
	"github.com/vidarandrebo/oncetree/protos/node"
	nmprotos "github.com/vidarandrebo/oncetree/protos/nodemanager"
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
	//	if logFile == io.Discard {
	//		logWriter = io.Discard
	//	}
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
func (n *Node) Run(knownAddr string, readyCallBack func()) {
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
	n.logger.Debug("RPC Crash")
	n.Stop("crash RPC")
	return &emptypb.Empty{}, nil
}
