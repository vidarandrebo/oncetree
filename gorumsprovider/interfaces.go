package gorumsprovider

import (
	"github.com/vidarandrebo/oncetree/protos/failuredetector"
	kvsprotos "github.com/vidarandrebo/oncetree/protos/keyvaluestorage"
)

type FDConfigProvider interface {
	FailureDetectorConfig() (*failuredetector.Configuration, bool)
}

type StorageConfigProvider interface {
	StorageConfig() (*kvsprotos.Configuration, bool, int)
	Reconnect(epoch int)
}
