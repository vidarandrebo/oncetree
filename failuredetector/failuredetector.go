package failuredetector

import (
	"context"
	"log"
	"reflect"
	"sync"
	"time"

	"github.com/vidarandrebo/oncetree/consts"
	"github.com/vidarandrebo/oncetree/eventbus"

	"github.com/vidarandrebo/oncetree/concurrent/hashset"
	"github.com/vidarandrebo/oncetree/concurrent/maps"
	"github.com/vidarandrebo/oncetree/nodemanager"

	"github.com/relab/gorums"

	fdprotos "github.com/vidarandrebo/oncetree/protos/failuredetector"
)

type FailureDetector struct {
	id            string
	nodes         *hashset.ConcurrentHashSet[string]
	alive         *maps.ConcurrentIntegerMap[string]
	suspected     *hashset.ConcurrentHashSet[string]
	logger        *log.Logger
	nodeManager   *nodemanager.NodeManager
	gorumsConfig  *fdprotos.Configuration
	gorumsManager *fdprotos.Manager
	eventBus      *eventbus.EventBus
}

func New(id string, logger *log.Logger, nodeManager *nodemanager.NodeManager, gorumsManager *fdprotos.Manager, eventBus *eventbus.EventBus) *FailureDetector {
	fd := &FailureDetector{
		id:            id,
		nodes:         hashset.New[string](),
		alive:         maps.NewConcurrentIntegerMap[string](),
		suspected:     hashset.New[string](),
		nodeManager:   nodeManager,
		gorumsManager: gorumsManager,
		eventBus:      eventBus,
		logger:        logger,
	}
	eventBus.RegisterHandler(
		reflect.TypeOf(nodemanager.NeighbourAddedEvent{}),
		func(e any) {
			if _, ok := e.(nodemanager.NeighbourAddedEvent); ok {
				fd.SetNodesFromManager()
				logger.Printf("fd tracking nodes %v", fd.nodes)
			}
		},
	)
	return fd
}

func (fd *FailureDetector) SetNodesFromManager() error {
	fd.nodes.Clear()
	fd.suspected.Clear()
	fd.alive.Clear()
	for _, neighbour := range fd.nodeManager.Neighbours() {
		fd.nodes.Add(neighbour.Key)
	}
	cfg, err := fd.gorumsManager.NewConfiguration(
		&qspec{
			numNodes: fd.nodes.Len(),
		},
		gorums.WithNodeMap(
			fd.nodeManager.GorumsNeighbourMap(),
		),
	)
	if err != nil {
		return err
	}
	fd.gorumsConfig = cfg
	return nil
}

func (fd *FailureDetector) DeregisterNode(nodeID string) {
	fd.nodes.Delete(nodeID)
	fd.alive.Delete(nodeID)
	fd.suspected.Delete(nodeID)
}

func (fd *FailureDetector) Suspect(nodeID string) {
	fd.eventBus.Push(NewNodeFailedEvent(nodeID))
}

func (fd *FailureDetector) Run(ctx context.Context, wg *sync.WaitGroup) {
mainLoop:
	for {
		select {
		case <-time.After(time.Duration(consts.FailureDetectorInterval) * time.Second):
			fd.timeout()
		case <-ctx.Done():
			break mainLoop
		case <-time.After(consts.HeartbeatSendInterval):
			fd.sendHeartbeat()
		}
	}
	fd.logger.Println("Exiting failure detector")
	wg.Done()
}

func (fd *FailureDetector) timeout() {
	for _, nodeID := range fd.nodes.Values() {
		if !fd.alive.Contains(nodeID) && !fd.suspected.Contains(nodeID) {
			fd.suspected.Add(nodeID)
			fd.logger.Printf("suspect node %v", nodeID)
			fd.Suspect(nodeID)
		}
	}

	fd.alive.Clear()
}

func (fd *FailureDetector) Heartbeat(ctx gorums.ServerCtx, request *fdprotos.HeartbeatMessage) {
	// fd.logger.Printf("received hb from %s", request.GetNodeID())
	fd.alive.Increment(request.GetNodeID(), 1)
}

func (fd *FailureDetector) sendHeartbeat() {
	go func() {
		if fd.gorumsConfig == nil {
			return
		}
		msg := fdprotos.HeartbeatMessage{NodeID: fd.id}
		ctx, cancel := context.WithTimeout(context.Background(), consts.RPCContextTimeout)
		defer cancel()
		fd.gorumsConfig.Heartbeat(ctx, &msg)
	}()
}
