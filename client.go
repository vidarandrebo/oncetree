package oncetree

import (
	"context"
	"fmt"
	"github.com/relab/gorums"
	"github.com/vidarandrebo/oncetree/protos"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"time"
)

type Client struct {
	manager *protos.Manager
	config  *protos.Configuration
	id      string
}

func NewClient(id string, nodes []string) *Client {
	manager := protos.NewManager(
		gorums.WithDialTimeout(1*time.Second),
		gorums.WithGrpcDialOptions(
			grpc.WithBlock(),
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		),
	)
	cfg, err := manager.NewConfiguration(&QSpec{numNodes: len(nodes)}, gorums.WithNodeList(nodes))
	if err != nil {
		log.Fatalln("failed to create gorums client gorumsConfig")
	}
	client := Client{
		id:      id,
		config:  cfg,
		manager: manager,
	}
	return &client
}
func (c *Client) Run() {
	_, err := c.config.Nodes()[0].Write(context.Background(), &protos.WriteRequest{
		Key:   99,
		Value: 100})
	if err != nil {
		fmt.Println(err)
	}
	response, err := c.config.Nodes()[0].Read(context.Background(), &protos.ReadRequest{Key: 99})
	if err == nil {
		fmt.Println(response)
	} else {
		fmt.Println(err)
	}
}
