package storage

import (
	"log/slog"
	"sync"

	"github.com/vidarandrebo/oncetree/gorumsprovider"
)

type GossipSender struct {
	logger         *slog.Logger
	workers        map[string]map[string]*GossipWorker // origin -> target -> Worker
	gorumsID       func(nodeID string) (uint32, bool)  // function to lookup gorumsID
	configProvider gorumsprovider.StorageConfigProvider
	mut            sync.Mutex
}

func NewGossipSender(logger *slog.Logger, configProvider gorumsprovider.StorageConfigProvider, gorumsIDFunc func(nodeID string) (uint32, bool)) *GossipSender {
	gs := GossipSender{
		logger:         logger.With(slog.Group("node", slog.String("module", "GossipSender"))),
		workers:        make(map[string]map[string]*GossipWorker),
		gorumsID:       gorumsIDFunc,
		configProvider: configProvider,
	}
	return &gs
}

func (gs *GossipSender) Enqueue(message PerNodeGossip, origin string, target string) {
	gs.mut.Lock()
	gs.ensureWorkerExists(origin, target)
	worker, ok := gs.workers[origin][target]
	gs.mut.Unlock()
	if !ok {
		panic("worker not found: " + origin)
	}
	worker.pending <- message
}

func (gs *GossipSender) ensureWorkerExists(origin string, target string) {
	if _, ok := gs.workers[origin]; !ok {
		gs.workers[origin] = make(map[string]*GossipWorker)
	}
	if _, ok := gs.workers[origin][target]; !ok {
		gs.workers[origin][target] = NewGossipWorker(gs.logger, gs.configProvider, gs.gorumsID, target)
	}
}
