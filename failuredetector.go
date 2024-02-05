package oncetree

import (
	"log"
	"sync"
	"time"
)

// FailureDetector is an implementation of EventuallyPerfectFailureDetector from
// "Introduction to Reliable and Secure Distributed Programming, second edition"
type FailureDetector struct {
	nodes     HashSet[string]
	alive     HashSet[string]
	suspected HashSet[string]
	delay     int
	logger    *log.Logger
	mut       sync.Mutex
}

func NewFailureDetector(logger *log.Logger) *FailureDetector {
	return &FailureDetector{
		nodes:     NewHashSet[string](),
		alive:     NewHashSet[string](),
		suspected: NewHashSet[string](),
		delay:     1,
		logger:    logger,
	}
}

func (fd *FailureDetector) RegisterHeartbeat(nodeID string) {
	fd.mut.Lock()
	fd.alive.Add(nodeID)
	fd.mut.Unlock()
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
	if len(fd.alive.Intersection(fd.suspected)) > 0 {
		fd.delay *= 2
		fd.logger.Printf("increase fd delay to %d s", fd.delay)
	}

	for nodeID := range fd.nodes {
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

	fd.mut.Lock()
	fd.alive = NewHashSet[string]()
	fd.mut.Unlock()
}
