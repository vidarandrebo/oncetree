package failuredetector

import (
	"context"
	"log/slog"
	"reflect"
	"sync"
	"time"

	"github.com/vidarandrebo/oncetree/storage/sevents"

	"github.com/vidarandrebo/oncetree/nodemanager/nmenums"

	"github.com/vidarandrebo/oncetree/failuredetector/fdevents"
	"github.com/vidarandrebo/oncetree/nodemanager/nmevents"

	"github.com/vidarandrebo/oncetree/gorumsprovider"

	"github.com/vidarandrebo/oncetree/consts"
	"github.com/vidarandrebo/oncetree/eventbus"

	"github.com/vidarandrebo/oncetree/common/concurrentmap"
	"github.com/vidarandrebo/oncetree/common/hashset"
	"github.com/vidarandrebo/oncetree/nodemanager"

	fdprotos "github.com/vidarandrebo/oncetree/protos/failuredetector"
)

type FailureDetector struct {
	id             string
	nodes          *hashset.ConcurrentHashSet[string]
	strikes        *concurrentmap.ConcurrentIntegerMap[string]
	suspected      *hashset.ConcurrentHashSet[string]
	logger         *slog.Logger
	nodeManager    *nodemanager.NodeManager
	eventBus       *eventbus.EventBus
	configProvider gorumsprovider.FDConfigProvider
	mut            sync.Mutex
}

func New(id string, logger *slog.Logger, nodeManager *nodemanager.NodeManager, eventBus *eventbus.EventBus, configProvider gorumsprovider.FDConfigProvider) *FailureDetector {
	fd := &FailureDetector{
		id:             id,
		nodes:          hashset.NewConcurrentHashSet[string](),
		strikes:        concurrentmap.NewConcurrentIntegerMap[string](),
		suspected:      hashset.NewConcurrentHashSet[string](),
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
	eventBus.RegisterHandler(
		reflect.TypeOf(sevents.GossipFailedEvent{}),
		func(e any) {
			if event, ok := e.(sevents.GossipFailedEvent); ok {
				fd.mut.Lock()
				if fd.nodes.Contains(event.NodeID) {
					fd.strikes.Increment(event.NodeID, 1)
				}
				fd.mut.Unlock()
			}
		})
	return fd
}

func (fd *FailureDetector) SetNodesFromManager() {
	fd.mut.Lock()
	defer fd.mut.Unlock()
	fd.nodes.Clear()
	fd.suspected.Clear()
	fd.strikes.Clear()
	for _, neighbour := range fd.nodeManager.Neighbours() {
		if (neighbour.Value.Role == nmenums.Parent) || (neighbour.Value.Role == nmenums.Child) {
			fd.nodes.Add(neighbour.Key)
		}
	}
}

func (fd *FailureDetector) Suspect(nodeID string) {
	fd.eventBus.PushEvent(fdevents.NewNodeFailedEvent(nodeID))
}

func (fd *FailureDetector) Run(ctx context.Context, wg *sync.WaitGroup) {
mainLoop:
	for {
		select {
		case <-time.After(consts.HeartbeatSendInterval):
			go fd.sendHeartbeat()
			fd.timeout()
		case <-ctx.Done():
			break mainLoop
		}
	}
	fd.logger.Debug("break loop")
	wg.Done()
}

func (fd *FailureDetector) timeout() {
	fd.mut.Lock()
	defer fd.mut.Unlock()
	for _, node := range fd.nodes.Values() {
		fd.strikes.Increment(node, 1)
	}
	for _, strike := range fd.strikes.Entries() {
		if strike.Value > consts.FailureDetectorStrikes && !fd.suspected.Contains(strike.Key) {
			fd.suspected.Add(strike.Key)
			fd.Suspect(strike.Key)
		}
	}
}

func (fd *FailureDetector) sendHeartbeat() {
	gorumsConfig, ok := fd.configProvider.FailureDetectorConfig()
	if !ok {
		fd.logger.Error("failed to retrieve config",
			"fn", "fd.sendHeartbeat")
		return
	}
	msg := fdprotos.HeartbeatMessage{NodeID: fd.id}
	ctx := context.Background()
	gorumsConfig.Heartbeat(ctx, &msg)
	if err := ctx.Err(); err != nil {
		fd.logger.Error("error when sending heartbeat",
			slog.Any("err", err))
	}
}
