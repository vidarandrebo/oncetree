package failuredetector

import (
	"context"
	"log/slog"
	"reflect"
	"sync"
	"time"

	"github.com/vidarandrebo/oncetree/nodemanager/nmenums"

	"github.com/vidarandrebo/oncetree/failuredetector/fdevents"
	"github.com/vidarandrebo/oncetree/nodemanager/nmevents"

	"github.com/vidarandrebo/oncetree/gorumsprovider"

	"github.com/vidarandrebo/oncetree/consts"
	"github.com/vidarandrebo/oncetree/eventbus"

	"github.com/vidarandrebo/oncetree/concurrent/hashset"
	"github.com/vidarandrebo/oncetree/concurrent/maps"
	"github.com/vidarandrebo/oncetree/nodemanager"

	fdprotos "github.com/vidarandrebo/oncetree/protos/failuredetector"
)

type FailureDetector struct {
	id             string
	nodes          *hashset.ConcurrentHashSet[string]
	strikes        *maps.ConcurrentIntegerMap[string]
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
		nodes:          hashset.New[string](),
		strikes:        maps.NewConcurrentIntegerMap[string](),
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
	nodeIDs := make([]string, 0)
	for _, gorumsID := range gorumsConfig.NodeIDs() {
		id, ok := fd.nodeManager.NodeID(gorumsID)
		if !ok {
			panic("jh")
		}
		nodeIDs = append(nodeIDs, id)
	}
	fd.logger.Debug("sending heartbeat",
		slog.Any("nodeIDs", nodeIDs))
	msg := fdprotos.HeartbeatMessage{NodeID: fd.id}
	ctx, cancel := context.WithTimeout(context.Background(), consts.RPCContextTimeout)
	defer cancel()
	gorumsConfig.Heartbeat(ctx, &msg)
	if err := ctx.Err(); err != nil {
		fd.logger.Error("error when sending heartbeat",
			slog.Any("err", err))
	}
}
