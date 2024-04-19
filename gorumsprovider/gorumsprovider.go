package gorumsprovider

import (
	"log/slog"
	"sync"

	fdprotos "github.com/vidarandrebo/oncetree/protos/failuredetector"
	kvsprotos "github.com/vidarandrebo/oncetree/protos/keyvaluestorage"
	"github.com/vidarandrebo/oncetree/protos/node"
	nmprotos "github.com/vidarandrebo/oncetree/protos/nodemanager"
)

// GorumsProvider manages gorums managers and provides gorums configs
type GorumsProvider struct {
	managers       *managers
	configurations *configurations
	mut            sync.RWMutex
	logger         *slog.Logger
}

func New(logger *slog.Logger) *GorumsProvider {
	return &GorumsProvider{
		managers:       newGorumsManagers(),
		configurations: newConfigurations(),
		logger:         logger.With(slog.Group("node", slog.String("module", "gorumsprovider"))),
	}
}

// SetNodes updates all configurations to contain the input nodes
func (gp *GorumsProvider) SetNodes(nodes map[string]uint32) {
	gp.mut.Lock()
	defer gp.mut.Unlock()
	gp.setNodes(nodes)
}

func (gp *GorumsProvider) setNodes(nodes map[string]uint32) {
	if len(nodes) == 0 {
		gp.logger.Warn("no nodes in node-map, skipping config creation")
		return
	}
	var err error

	// NodeManager
	gp.configurations.nmConfig, err = gp.managers.newNodeManagerConfig(nodes)
	if err != nil {
		gp.logger.Error("failed to create nodemanager config", "err", err)
	} else {
		// gp.logger.Println("[GorumsProvider] - Created nodemanager config")
	}

	// Node
	gp.configurations.nodeConfig, err = gp.managers.newNodeConfig(nodes)
	if err != nil {
		gp.logger.Error("failed to create node config", "err", err)
	} else {
		// gp.logger.Println("[GorumsProvider] - Created node config")
	}

	// StorageService
	gp.configurations.kvsConfig, err = gp.managers.newKVSConfig(nodes)
	if err != nil {
		gp.logger.Error("failed to create storage config", "err", err)
	} else {
		// gp.logger.Println("[GorumsProvider] - Created storage config")
	}

	// FailureDetector
	gp.configurations.fdConfig, err = gp.managers.newFDConfig(nodes)
	if err != nil {
		gp.logger.Error("failed to create failuredetector config", "err", err)
	} else {
		gp.logger.Debug("created failuredetector config")
	}
}

// Reset deletes existing manager and configs, then it creates new ones
func (gp *GorumsProvider) Reset() {
	gp.mut.Lock()
	defer gp.mut.Unlock()
	gp.reset()
}

func (gp *GorumsProvider) reset() {
	gp.managers = gp.managers.recreate(gp.logger)
	gp.configurations = newConfigurations()
}

// ResetWithNewNodes deletes existing manager and configs, then it creates new ones with the provided nodes
func (gp *GorumsProvider) ResetWithNewNodes(nodes map[string]uint32) {
	gp.mut.Lock()
	defer gp.mut.Unlock()
	gp.reset()
	gp.setNodes(nodes)
}

func (gp *GorumsProvider) FailureDetectorConfig() (*fdprotos.Configuration, bool) {
	gp.mut.RLock()
	defer gp.mut.RUnlock()
	if gp.configurations.fdConfig == nil {
		return nil, false
	}
	return gp.configurations.fdConfig, true
}

func (gp *GorumsProvider) NodeManagerConfig() (*nmprotos.Configuration, bool) {
	gp.mut.RLock()
	defer gp.mut.RUnlock()
	if gp.configurations.nmConfig == nil {
		return nil, false
	}
	return gp.configurations.nmConfig, true
}

func (gp *GorumsProvider) CustomNodeManagerConfig(gorumsMap map[string]uint32) (*nmprotos.Configuration, error) {
	gp.mut.RLock()
	defer gp.mut.RUnlock()
	cfg, err := gp.managers.newNodeManagerConfig(gorumsMap)
	return cfg, err
}

func (gp *GorumsProvider) NodeConfig() (*node.Configuration, bool) {
	gp.mut.RLock()
	defer gp.mut.RUnlock()
	if gp.configurations.nodeConfig == nil {
		return nil, false
	}
	return gp.configurations.nodeConfig, true
}

func (gp *GorumsProvider) StorageConfig() (*kvsprotos.Configuration, bool) {
	gp.mut.RLock()
	defer gp.mut.RUnlock()
	if gp.configurations.kvsConfig == nil {
		return nil, false
	}
	return gp.configurations.kvsConfig, true
}
