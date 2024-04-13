package failuredetector

import (
	"context"
	"log/slog"
	"reflect"
	"sync"
	"time"

	"github.com/vidarandrebo/oncetree/failuredetector/fdevents"
	"github.com/vidarandrebo/oncetree/nodemanager/nmevents"

	"github.com/vidarandrebo/oncetree/gorumsprovider"

	"github.com/vidarandrebo/oncetree/consts"
	"github.com/vidarandrebo/oncetree/eventbus"

	"github.com/vidarandrebo/oncetree/concurrent/hashset"
	"github.com/vidarandrebo/oncetree/concurrent/maps"
	"github.com/vidarandrebo/oncetree/nodemanager"

	"github.com/relab/gorums"

	fdprotos "github.com/vidarandrebo/oncetree/protos/failuredetector"
)

type FailureDetector struct {
	id             string
	nodes          *hashset.ConcurrentHashSet[string]
	alive          *maps.ConcurrentIntegerMap[string]
	suspected      *hashset.ConcurrentHashSet[string]
	logger         *slog.Logger
	nodeManager    *nodemanager.NodeManager
	eventBus       *eventbus.EventBus
	configProvider gorumsprovider.FDConfigProvider
}

func New(id string, logger *slog.Logger, nodeManager *nodemanager.NodeManager, eventBus *eventbus.EventBus, configProvider gorumsprovider.FDConfigProvider) *FailureDetector {
	fd := &FailureDetector{
		id:             id,
		nodes:          hashset.New[string](),
		alive:          maps.NewConcurrentIntegerMap[string](),
		suspected:      hashset.New[string](),
		nodeManager:    nodeManager,
		eventBus:       eventBus,
		logger:         logger.With(slog.Group("node", slog.String("module", "failuredetector"))),
		configProvider: configProvider,
	}
	eventBus.RegisterHandler(
		reflect.TypeOf(nmevents.NeighbourAddedEvent{}),
		func(e any) {
			if _, ok := e.(nmevents.NeighbourAddedEvent); ok {
				fd.SetNodesFromManager()
				// logger.Printf("fd tracking nodes %v", fd.nodes)
			}
		},
	)
	return fd
}

func (fd *FailureDetector) SetNodesFromManager() {
	fd.nodes.Clear()
	fd.suspected.Clear()
	fd.alive.Clear()
	for _, neighbour := range fd.nodeManager.Neighbours() {
		if (neighbour.Value.Role == nodemanager.Parent) || (neighbour.Value.Role == nodemanager.Child) {
			fd.nodes.Add(neighbour.Key)
			fd.alive.Increment(neighbour.Key, 1)
		}
	}
}

func (fd *FailureDetector) DeregisterNode(nodeID string) {
	fd.nodes.Delete(nodeID)
	fd.alive.Delete(nodeID)
	fd.suspected.Delete(nodeID)
}

func (fd *FailureDetector) Suspect(nodeID string) {
	fd.eventBus.PushEvent(fdevents.NewNodeFailedEvent(nodeID))
}

func (fd *FailureDetector) Run(ctx context.Context, wg *sync.WaitGroup) {
	heartbeatTimeout := time.After(consts.FailureDetectionInterval)
	heartbeatTick := time.Tick(consts.HeartbeatSendInterval)

mainLoop:
	for {
		select {
		case <-heartbeatTimeout:
			fd.timeout()
			// timeout operations should not be queued,
			// so therefore the time.After is reset instead of using ticker
			heartbeatTimeout = time.After(consts.FailureDetectionInterval)
		case <-ctx.Done():
			break mainLoop
		case <-heartbeatTick:
			fd.sendHeartbeat()
		}
	}
	fd.logger.Debug("break loop")
	wg.Done()
}

func (fd *FailureDetector) timeout() {
	fd.logger.Debug("timeout")
	for _, nodeID := range fd.nodes.Values() {
		if !fd.alive.Contains(nodeID) && !fd.suspected.Contains(nodeID) {
			fd.suspected.Add(nodeID)
			fd.logger.Warn(
				"suspect",
				slog.String("node", nodeID))
			fd.Suspect(nodeID)
		}
	}

	fd.alive.Clear()
}

func (fd *FailureDetector) Heartbeat(ctx gorums.ServerCtx, request *fdprotos.HeartbeatMessage) {
	fd.logger.Debug(
		"RPC Heartbeat",
		slog.String("id", request.GetNodeID()),
	)
	fd.alive.Increment(request.GetNodeID(), 1)
	if fd.suspected.Contains(request.GetNodeID()) {
		fd.logger.Error("received heartbeat from suspected node", slog.String("id", request.GetNodeID()))
		panic("heartbeat problem")
	}
}

func (fd *FailureDetector) sendHeartbeat() {
	fd.logger.Debug("send heartbeat")
	gorumsConfig, ok := fd.configProvider.FailureDetectorConfig()
	if !ok {
		fd.logger.Error("failed to retrieve config",
			"fn", "fd.sendHeartbeat")
		return
	}
	msg := fdprotos.HeartbeatMessage{NodeID: fd.id}
	ctx, cancel := context.WithTimeout(context.Background(), consts.RPCContextTimeout)
	defer cancel()
	gorumsConfig.Heartbeat(ctx, &msg)
	if err := ctx.Err(); err != nil {
		fd.logger.Error("error when sending heartbeat",
			slog.Any("err", err))
	}
}
