package oncetree

import (
	"context"
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
	cfg := c.gorumsProvider.NodeConfig()
	node, exists := cfg.Node(0)
	if !exists {
		panic("node does not exists")
	}
	node.Crash(context.Background(), &emptypb.Empty{})
}
