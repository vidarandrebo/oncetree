package oncetree

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/relab/gorums"
	"github.com/vidarandrebo/oncetree/protos"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Client struct {
	manager *protos.Manager
	config  *protos.Configuration
}

func NewClient(nodes []string) *Client {
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
		config:  cfg,
		manager: manager,
	}
	return &client
}

func (c *Client) Run() {
	for _, node := range c.config.Nodes() {
		_, err := node.Write(context.Background(), &protos.WriteRequest{
			Key:   20,
			Value: 10,
		})
		if err != nil {
			fmt.Println(err)
		}
	}
	time.Sleep(5 * time.Second)

	for _, node := range c.config.Nodes() {
		_, err := node.PrintState(context.Background(), &emptypb.Empty{})
		if err != nil {
			log.Fatalln(err)
		}
	}
}
