package gorumsprovider

import (
	fdprotos "github.com/vidarandrebo/oncetree/protos/failuredetector"
	kvsprotos "github.com/vidarandrebo/oncetree/protos/keyvaluestorage"
	"github.com/vidarandrebo/oncetree/protos/node"
	nmprotos "github.com/vidarandrebo/oncetree/protos/nodemanager"
)

type configurations struct {
	fdConfig   *fdprotos.Configuration
	nmConfig   *nmprotos.Configuration
	kvsConfig  *kvsprotos.Configuration
	nodeConfig *node.Configuration
}

func newConfigurations() *configurations {
	return &configurations{}
}
