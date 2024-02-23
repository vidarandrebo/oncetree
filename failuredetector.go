package oncetree

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/relab/gorums"

	fdprotos "github.com/vidarandrebo/oncetree/protos/failuredetectorprotos"
)

// FailureDetector is an implementation of EventuallyPerfectFailureDetector from
// "Introduction to Reliable and Secure Distributed Programming, second edition"
type FailureDetector struct {
	id            string
	nodes         *ConcurrentHashSet[string]
	alive         *ConcurrentIntegerMap[string]
	suspected     *ConcurrentHashSet[string]
	delay         int
	logger        *log.Logger
	subscribers   []chan<- string
	configuration *fdprotos.Configuration
	manager       *fdprotos.Manager
}

func NewFailureDetector(id string, logger *log.Logger) *FailureDetector {
	return &FailureDetector{
		id:        id,
		nodes:     NewConcurrentHashSet[string](),
		alive:     NewConcurrentIntegerMap[string](),
		suspected: NewConcurrentHashSet[string](),
		delay:     5,
		logger:    logger,
	}
}

func (fd *FailureDetector) RegisterHeartbeat(nodeID string) {
	fd.alive.Increment(nodeID, 1)
}

func (fd *FailureDetector) RegisterNode(nodeID string) {
	fd.nodes.Add(nodeID)
}

func (fd *FailureDetector) RegisterNodes(nodeIDs []string) {
	for _, nodeID := range nodeIDs {
		fd.nodes.Add(nodeID)
	}
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
	go func() {
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
	}()
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
	fd.RegisterHeartbeat(request.GetNodeID())
}

func (fd *FailureDetector) sendHeartbeat() {
	go func() {
		if fd.configuration == nil {
			return
		}
		msg := fdprotos.HeartbeatMessage{NodeID: fd.id}
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		fd.configuration.Heartbeat(ctx, &msg)
	}()
}
