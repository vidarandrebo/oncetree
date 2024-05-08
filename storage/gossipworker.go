package storage

import (
	"context"
	"github.com/vidarandrebo/oncetree/gorumsprovider"
	kvsprotos "github.com/vidarandrebo/oncetree/protos/keyvaluestorage"
)

type GossipWorker struct {
	configProvider gorumsprovider.StorageConfigProvider
	pending        chan PerNodeGossip
	target         string
	gorumsID       func(nodeID string) (uint32, bool)
}

func NewGossipWorker(configProvider gorumsprovider.StorageConfigProvider, gorumsIDFunc func(nodeID string) (uint32, bool), target string) *GossipWorker {
	worker := &GossipWorker{
		pending:        make(chan PerNodeGossip),
		configProvider: configProvider,
		target:         target,
		gorumsID:       gorumsIDFunc,
	}
	go worker.Run()
	return worker
}

func (gw *GossipWorker) Run() {
	for {
		select {
		case gossip := <-gw.pending:
			gw.sendGossip(gossip)
		}
	}
}

func (gw *GossipWorker) sendGossip(gossip PerNodeGossip) {
	cfg, ok, _ := gw.configProvider.StorageConfig()
	if !ok {
		return
	}
	gorumsID, ok := gw.gorumsID(gw.target)
	if !ok {
		return
	}
	node, ok := cfg.Node(gorumsID)
	if !ok {
		return
	}
	ctx := context.Background()
	message := kvsprotos.GossipMessage{
		NodeID:         gossip.NodeID,
		Key:            gossip.Key,
		AggValue:       gossip.AggValue,
		AggTimestamp:   gossip.AggTimestamp,
		LocalValue:     gossip.LocalValue,
		LocalTimestamp: gossip.LocalTimestamp,
	}
	node.Gossip(ctx, &message)
}
