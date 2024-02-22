package oncetree

import (
	"time"

	"github.com/relab/gorums"
	fdprotos "github.com/vidarandrebo/oncetree/protos/failuredetectorprotos"
	kvsprotos "github.com/vidarandrebo/oncetree/protos/keyvaluestorageprotos"
	"github.com/vidarandrebo/oncetree/protos/nodeprotos"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type GorumsManagers struct {
	fdManager   *fdprotos.Manager
	kvsManager  *kvsprotos.Manager
	nodeManager *nodeprotos.Manager
}

func CreateGorumsManagers() *GorumsManagers {
	fdManager := fdprotos.NewManager(
		gorums.WithDialTimeout(1*time.Second),
		gorums.WithGrpcDialOptions(
			grpc.WithBlock(),
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		),
	)
	kvsManager := kvsprotos.NewManager(
		gorums.WithDialTimeout(1*time.Second),
		gorums.WithGrpcDialOptions(
			grpc.WithBlock(),
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		),
	)
	nodeManager := nodeprotos.NewManager(
		gorums.WithDialTimeout(1*time.Second),
		gorums.WithGrpcDialOptions(
			grpc.WithBlock(),
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		),
	)
	return &GorumsManagers{
		fdManager:   fdManager,
		kvsManager:  kvsManager,
		nodeManager: nodeManager,
	}
}
