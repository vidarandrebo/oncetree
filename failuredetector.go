package oncetree

import (
	"context"
	"log"
	"sync"
	"time"
)

// FailureDetector is an implementation of EventuallyPerfectFailureDetector from
// "Introduction to Reliable and Secure Distributed Programming, second edition"
type FailureDetector struct {
	nodes       *ConcurrentHashSet[string]
	alive       *ConcurrentIntegerMap[string]
	suspected   *ConcurrentHashSet[string]
	delay       int
	logger      *log.Logger
	subscribers []chan<- string
}

func NewFailureDetector(logger *log.Logger) *FailureDetector {
	return &FailureDetector{
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
