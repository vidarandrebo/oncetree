package oncetree

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/relab/gorums"

	fdprotos "github.com/vidarandrebo/oncetree/protos/failuredetector"
)

type FailureDetector struct {
	id            string
	nodes         *ConcurrentHashSet[string]
	alive         *ConcurrentIntegerMap[string]
	suspected     *ConcurrentHashSet[string]
	delay         int
	logger        *log.Logger
	subscribers   []chan<- string
	nodeManager   *NodeManager
	gorumsConfig  *fdprotos.Configuration
	gorumsManager *fdprotos.Manager
}

func NewFailureDetector(id string, logger *log.Logger, nodeManager *NodeManager, gorumsManager *fdprotos.Manager) *FailureDetector {
	return &FailureDetector{
		id:            id,
		nodes:         NewConcurrentHashSet[string](),
		alive:         NewConcurrentIntegerMap[string](),
		suspected:     NewConcurrentHashSet[string](),
		nodeManager:   nodeManager,
		gorumsManager: gorumsManager,
		delay:         5,
		logger:        logger,
		subscribers:   make([]chan<- string, 0),
	}
}

func (fd *FailureDetector) SetNodesFromManager() error {
	for _, neighbour := range fd.nodeManager.GetNeighbours() {
		fd.nodes.Add(neighbour.Key)
	}
	cfg, err := fd.gorumsManager.NewConfiguration(
		&FDQSpec{
			NumNodes: fd.nodes.Len(),
		},
		gorums.WithNodeMap(
			fd.nodeManager.GetGorumsNeighbourMap(),
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

func (fd *FailureDetector) Subscribe() <-chan string {
	channel := make(chan string)
	fd.subscribers = append(fd.subscribers, channel)
	return channel
}

func (fd *FailureDetector) Suspect(nodeID string) {
	for _, c := range fd.subscribers {
		c <- nodeID
	}
}

func (fd *FailureDetector) Run(ctx context.Context, wg *sync.WaitGroup) {
mainLoop:
	for {
		select {
		case <-time.After(time.Duration(fd.delay) * time.Second):
			fd.timeout()
		case <-ctx.Done():
			break mainLoop
		case <-time.After(time.Second * 1):
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
	fd.alive.Increment(request.GetNodeID(), 1)
}

func (fd *FailureDetector) sendHeartbeat() {
	go func() {
		if fd.gorumsConfig == nil {
			return
		}
		msg := fdprotos.HeartbeatMessage{NodeID: fd.id}
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		fd.gorumsConfig.Heartbeat(ctx, &msg)
	}()
}
