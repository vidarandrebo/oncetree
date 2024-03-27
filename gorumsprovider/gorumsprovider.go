package gorumsprovider

import (
	"log"
	"sync"

	"github.com/relab/gorums"
	"github.com/vidarandrebo/oncetree/failuredetector/fdqspec"
	"github.com/vidarandrebo/oncetree/nodemanager/nmqspec"
	fdprotos "github.com/vidarandrebo/oncetree/protos/failuredetector"
	kvsprotos "github.com/vidarandrebo/oncetree/protos/keyvaluestorage"
	"github.com/vidarandrebo/oncetree/protos/node"
	nmprotos "github.com/vidarandrebo/oncetree/protos/nodemanager"
	"github.com/vidarandrebo/oncetree/storage/sqspec"
)

// GorumsProvider manages gorums managers and provides gorums configs
type GorumsProvider struct {
	managers       *managers
	configurations *configurations
	mut            sync.RWMutex
	logger         *log.Logger
}

func New(logger *log.Logger) *GorumsProvider {
	return &GorumsProvider{
		managers:       newGorumsManagers(),
		configurations: newConfigurations(),
		logger:         logger,
	}
}

// SetNodes updates all configurations to contain the input nodes
func (gp *GorumsProvider) SetNodes(nodes map[string]uint32) {
	gp.mut.Lock()
	defer gp.mut.Unlock()
	var err error

	// NodeManager
	gp.configurations.nmConfig, err = gp.managers.nmManager.NewConfiguration(
		&nmqspec.QSpec{
			NumNodes: len(nodes),
		},
		gorums.WithNodeMap(nodes),
	)

	if err != nil {
		gp.logger.Println(err)
		gp.logger.Println("[GorumsProvider] - Failed to create nodemanager config")
	} else {
		// gp.logger.Println("[GorumsProvider] - Created nodemanager config")
	}

	// Node
	gp.configurations.nodeConfig, err = gp.managers.nodeManager.NewConfiguration(
		gorums.WithNodeMap(nodes),
	)
	if err != nil {
		gp.logger.Println(err)
		gp.logger.Println("[GorumsProvider] - Failed to create node config")
	} else {
		// gp.logger.Println("[GorumsProvider] - Created node config")
	}

	// StorageService
	gp.configurations.kvsConfig, err = gp.managers.kvsManager.NewConfiguration(
		&sqspec.QSpec{
			NumNodes: len(nodes),
		},
		gorums.WithNodeMap(nodes),
	)
	if err != nil {
		gp.logger.Println(err)
		gp.logger.Println("[GorumsProvider] - Failed to create storage config")
	} else {
		// gp.logger.Println("[GorumsProvider] - Created storage config")
	}

	// FailureDetector
	gp.configurations.fdConfig, err = gp.managers.fdManager.NewConfiguration(
		&fdqspec.QSpec{
			NumNodes: len(nodes),
		},
		gorums.WithNodeMap(nodes),
	)
	if err != nil {
		gp.logger.Println(err)
		gp.logger.Println("[GorumsProvider] - Failed to create failuredetector config")
	} else {
		// gp.logger.Println("[GorumsProvider] - Created failuredetector config")
	}
}

func (gp *GorumsProvider) FailureDetectorConfig() *fdprotos.Configuration {
	gp.mut.RLock()
	defer gp.mut.RUnlock()
	return gp.configurations.fdConfig
}

func (gp *GorumsProvider) NodeManagerConfig() *nmprotos.Configuration {
	gp.mut.RLock()
	defer gp.mut.RUnlock()
	return gp.configurations.nmConfig
}

func (gp *GorumsProvider) NodeConfig() *node.Configuration {
	gp.mut.RLock()
	defer gp.mut.RUnlock()
	return gp.configurations.nodeConfig
}

func (gp *GorumsProvider) StorageConfig() *kvsprotos.Configuration {
	gp.mut.RLock()
	defer gp.mut.RUnlock()
	return gp.configurations.kvsConfig
}

// Reset deletes existing manager and configs, then it creates new ones
func (gp *GorumsProvider) Reset() {
	gp.managers.recreate()
	gp.configurations = newConfigurations()
}
