package oncetree

import (
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
}

func NewClient(nodes []string) *Client {
	client := Client{}
	client.manager = protos.NewManager(
		gorums.WithDialTimeout(1*time.Second),
		gorums.WithGrpcDialOptions(
			grpc.WithBlock(),
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		),
	)
	cfg, err := client.manager.NewConfiguration(QSpec{numNodes: len(nodes)}, gorums.WithNodeList(nodes))
	if err != nil {
		log.Fatalln("failed to create gorums client config")
	}
	client.config = cfg
	return &client
}
func (c *Client) Run() {
	//response, err := c.config.Nodes()[0].Write(context.Background(), &protos.WriteRequest{Key: 99, Value: 100})
}

type QSpec struct {
	numNodes int
}
