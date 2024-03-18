package gorumsprovider

import (
	"github.com/relab/gorums"
	"github.com/vidarandrebo/oncetree/consts"
	fdprotos "github.com/vidarandrebo/oncetree/protos/failuredetector"
	kvsprotos "github.com/vidarandrebo/oncetree/protos/keyvaluestorage"
	"github.com/vidarandrebo/oncetree/protos/node"
	nmprotos "github.com/vidarandrebo/oncetree/protos/nodemanager"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type managers struct {
	fdManager   *fdprotos.Manager
	kvsManager  *kvsprotos.Manager
	nmManager   *nmprotos.Manager
	nodeManager *node.Manager
}

func newGorumsManagers() *managers {
	opts := []gorums.ManagerOption{
		gorums.WithDialTimeout(consts.GorumsDialTimeout),
		gorums.WithGrpcDialOptions(
			grpc.WithBlock(),
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		),
	}
	fdManager := fdprotos.NewManager(
		opts...,
	)
	kvsManager := kvsprotos.NewManager(
		opts...,
	)
	nmManager := nmprotos.NewManager(
		opts...,
	)
	nodeManager := node.NewManager(
		opts...,
	)

	return &managers{
		fdManager:   fdManager,
		kvsManager:  kvsManager,
		nmManager:   nmManager,
		nodeManager: nodeManager,
	}
}

func (m *managers) recreate() *managers {
	m.fdManager.Close()
	m.nmManager.Close()
	m.kvsManager.Close()
	m.nodeManager.Close()
	return newGorumsManagers()
}
