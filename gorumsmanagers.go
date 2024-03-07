package oncetree

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

type GorumsManagers struct {
	fdManager   *fdprotos.Manager
	kvsManager  *kvsprotos.Manager
	nmManager   *nmprotos.Manager
	nodeManager *node.Manager
}

func NewGorumsManagers() *GorumsManagers {
	fdManager := fdprotos.NewManager(
		gorums.WithDialTimeout(consts.GorumsDialTimeout),
		gorums.WithGrpcDialOptions(
			grpc.WithBlock(),
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		),
	)
	kvsManager := kvsprotos.NewManager(
		gorums.WithDialTimeout(consts.GorumsDialTimeout),
		gorums.WithGrpcDialOptions(
			grpc.WithBlock(),
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		),
	)
	nmManager := nmprotos.NewManager(
		gorums.WithDialTimeout(consts.GorumsDialTimeout),
		gorums.WithGrpcDialOptions(
			grpc.WithBlock(),
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		),
	)
	nodeManager := node.NewManager(
		gorums.WithDialTimeout(consts.GorumsDialTimeout),
		gorums.WithGrpcDialOptions(
			grpc.WithBlock(),
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		),
	)

	return &GorumsManagers{
		fdManager:   fdManager,
		kvsManager:  kvsManager,
		nmManager:   nmManager,
		nodeManager: nodeManager,
	}
}
