package oncetree

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/vidarandrebo/oncetree/storage/sqspec"

	"github.com/relab/gorums"
	kvsprotos "github.com/vidarandrebo/oncetree/protos/keyvaluestorage"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Client struct {
	config  *kvsprotos.Configuration
	manager *kvsprotos.Manager
}

func NewClient(nodes []string) *Client {
	manager := kvsprotos.NewManager(
		gorums.WithDialTimeout(1*time.Second),
		gorums.WithGrpcDialOptions(
			grpc.WithBlock(),
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		),
	)
	cfg, err := manager.NewConfiguration(&sqspec.QSpec{NumNodes: len(nodes)}, gorums.WithNodeList(nodes))
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
		log.Printf("sending to %v", node.Address())
		_, err := node.Write(context.Background(), &kvsprotos.WriteRequest{
			Key:   20,
			Value: 10,
		})
		if err != nil {
			fmt.Println(err)
		}
	}
	log.Println("hello")
	time.Sleep(5 * time.Second)

	for _, node := range c.config.Nodes() {
		_, err := node.PrintState(context.Background(), &emptypb.Empty{})
		if err != nil {
			log.Fatalln(err)
		}
	}
}
