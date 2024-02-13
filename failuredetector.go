package oncetree

import (
	"log"
	"time"
)

// FailureDetector is an implementation of EventuallyPerfectFailureDetector from
// "Introduction to Reliable and Secure Distributed Programming, second edition"
type FailureDetector struct {
	nodes     *ConcurrentHashSet[string]
	alive     *ConcurrentHashSet[string]
	suspected *ConcurrentHashSet[string]
	delay     int
	logger    *log.Logger
}

func NewFailureDetector(logger *log.Logger) *FailureDetector {
	return &FailureDetector{
		nodes:     NewConcurrentHashSet[string](),
		alive:     NewConcurrentHashSet[string](),
		suspected: NewConcurrentHashSet[string](),
		delay:     1,
		logger:    logger,
	}
}

func (fd *FailureDetector) RegisterHeartbeat(nodeID string) {
	fd.alive.Add(nodeID)
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
	fd.nodes.Remove(nodeID)
	fd.alive.Remove(nodeID)
	fd.suspected.Remove(nodeID)
}

func (fd *FailureDetector) Run() {
	go func() {
		for {
			select {
			case <-time.After(time.Duration(fd.delay) * time.Second):
				fd.timeout()
			}
		}
	}()
}

func (fd *FailureDetector) timeout() {
	// double delay if nodes are in both alive and suspected sets
	if fd.alive.Intersection(fd.suspected).Len() > 0 {
		fd.delay *= 2
		fd.logger.Printf("increase fd delay to %d s", fd.delay)
	}

	for _, nodeID := range fd.nodes.Values() {
		if !fd.alive.Contains(nodeID) && !fd.suspected.Contains(nodeID) {
			fd.suspected.Add(nodeID)
			fd.logger.Printf("suspect node %v", nodeID)
			// trigger suspect operation
		} else if fd.alive.Contains(nodeID) && fd.suspected.Contains(nodeID) {
			fd.suspected.Remove(nodeID)
			fd.logger.Printf("restore node %v", nodeID)
			// trigger restore
		}
	}

	fd.alive.Clear()
}
