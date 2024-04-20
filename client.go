package oncetree

import (
	"context"
	"fmt"
	kvsprotos "github.com/vidarandrebo/oncetree/protos/keyvaluestorage"
	"log/slog"

	"github.com/vidarandrebo/oncetree/gorumsprovider"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Client struct {
	gorumsProvider *gorumsprovider.GorumsProvider
}

func NewClient(nodes map[string]uint32) *Client {
	provider := gorumsprovider.New(slog.Default())
	provider.SetNodes(nodes)
	return &Client{gorumsProvider: provider}
}

func (c *Client) Run() {
	nodeCfg, ok := c.gorumsProvider.NodeConfig()
	if !ok {
		return
	}
	storageCfg, ok := c.gorumsProvider.StorageConfig()
	if !ok {
		return
	}
	storageNode, exists := storageCfg.Node(1)
	if !exists {
		panic("node does not exists")
	}
	for i := 0; i < 10000; i++ {
		response, err := storageNode.Write(context.Background(), &kvsprotos.WriteRequest{
			Key:     int64((i % 100) + 1),
			Value:   int64(i + 2),
			WriteID: int64(i + 1),
		})
		if err != nil {
			fmt.Println(response)
			panic("write failed")
		}
	}
	node, exists := nodeCfg.Node(1)
	if !exists {
		panic("node does not exists")
	}
	node.Crash(context.Background(), &emptypb.Empty{})
}
