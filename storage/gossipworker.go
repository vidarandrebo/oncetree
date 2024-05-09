package storage

import (
	"context"
	"log/slog"

	"github.com/vidarandrebo/oncetree/consts"

	"github.com/vidarandrebo/oncetree/gorumsprovider"
	kvsprotos "github.com/vidarandrebo/oncetree/protos/keyvaluestorage"
)

type GossipWorker struct {
	logger         *slog.Logger
	configProvider gorumsprovider.StorageConfigProvider
	pending        chan PerNodeGossip
	target         string
	gorumsID       func(nodeID string) (uint32, bool)
	hasSkipped     bool
}

func NewGossipWorker(logger *slog.Logger, configProvider gorumsprovider.StorageConfigProvider, gorumsIDFunc func(nodeID string) (uint32, bool), target string) *GossipWorker {
	worker := &GossipWorker{
		logger:         logger.With(slog.Group("node", slog.String("module", "GossipWorker"))),
		pending:        make(chan PerNodeGossip, consts.GossipWorkerBuffSize),
		configProvider: configProvider,
		target:         target,
		gorumsID:       gorumsIDFunc,
		hasSkipped:     false,
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
		gw.logger.Error("no cfg, skipping gossip",
			slog.String("target", gw.target))
		return
	}
	gorumsID, ok := gw.gorumsID(gw.target)
	if !ok {
		if !gw.hasSkipped {
			gw.logger.Error("node not in manager, skipping",
				slog.String("target", gw.target))
		}
		gw.hasSkipped = true
		return
	}
	node, ok := cfg.Node(gorumsID)
	if !ok {
		gw.logger.Error("node not in config, skipping",
			slog.String("target", gw.target))
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
	_, err := node.Gossip(ctx, &message)
	if err != nil {
		gw.logger.Error("sending of gossip message failed",
			slog.Any("err", err),
			slog.Int64("key", gossip.Key),
			slog.String("target", gw.target),
			slog.Int64("value", gossip.AggValue))
	}
}
