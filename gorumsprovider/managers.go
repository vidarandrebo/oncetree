package gorumsprovider

import (
	"log/slog"
	"time"

	"github.com/relab/gorums"
	"github.com/vidarandrebo/oncetree/consts"
	"github.com/vidarandrebo/oncetree/failuredetector/fdqspec"
	"github.com/vidarandrebo/oncetree/nodemanager/nmqspec"
	fdprotos "github.com/vidarandrebo/oncetree/protos/failuredetector"
	kvsprotos "github.com/vidarandrebo/oncetree/protos/keyvaluestorage"
	"github.com/vidarandrebo/oncetree/protos/node"
	nmprotos "github.com/vidarandrebo/oncetree/protos/nodemanager"
	"github.com/vidarandrebo/oncetree/storage/sqspec"
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

// recreate schedules a closing of gorums-managers in managers object and returns a new *managers reference
func (m *managers) recreate(logger *slog.Logger) *managers {
	go m.closeManagers(logger)
	logger.Info("gorums manager disposal queued")
	return newGorumsManagers()
}

// closeManagers closes all stored manager and its connections, should be invoked as a goroutine
// Managers are closed after a delay so that goroutines using a configuration are not interrupted until after the delay.
// By this time, no goroutines should be using the old configurations.
func (m *managers) closeManagers(logger *slog.Logger) {
	time.Sleep(consts.CloseMgrDelay)
	m.fdManager.Close()
	m.nmManager.Close()
	m.kvsManager.Close()
	m.nodeManager.Close()
	logger.Debug("gorums manager's connections disposed")
}

func (m *managers) newFDConfig(nodes map[string]uint32) (*fdprotos.Configuration, error) {
	cfg, err := m.fdManager.NewConfiguration(
		&fdqspec.QSpec{
			NumNodes: len(nodes),
		},
		gorums.WithNodeMap(nodes),
	)
	return cfg, err
}

func (m *managers) newKVSConfig(nodes map[string]uint32) (*kvsprotos.Configuration, error) {
	cfg, err := m.kvsManager.NewConfiguration(
		&sqspec.QSpec{
			NumNodes: len(nodes),
		},
		gorums.WithNodeMap(nodes),
	)
	return cfg, err
}

func (m *managers) newNodeConfig(nodes map[string]uint32) (*node.Configuration, error) {
	cfg, err := m.nodeManager.NewConfiguration(
		gorums.WithNodeMap(nodes),
	)
	return cfg, err
}

func (m *managers) newNodeManagerConfig(nodes map[string]uint32) (*nmprotos.Configuration, error) {
	cfg, err := m.nmManager.NewConfiguration(
		&nmqspec.QSpec{
			NumNodes: len(nodes),
		},
		gorums.WithNodeMap(nodes),
	)
	return cfg, err
}
