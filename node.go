package oncetree

import (
	"github.com/relab/gorums"
	"github.com/vidarandrebo/oncetree/protos"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/emptypb"
	"log"
	"net"
	"time"
)

type Node struct {
	rpcAddr         string
	keyValueStorage *KeyValueStorage
	gorumsServer    *gorums.Server
	id              string
	parent          string
	children        []string
	manager         *protos.Manager
	config          *protos.Configuration
}

func NewNode(id string, rpcAddr string) *Node {
	manager := protos.NewManager(
		gorums.WithDialTimeout(1*time.Second),
		gorums.WithGrpcDialOptions(
			grpc.WithBlock(),
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		),
	)
	return &Node{
		rpcAddr:         rpcAddr,
		keyValueStorage: NewKeyValueStorage(),
		id:              id,
		parent:          "",
		children:        make([]string, 0),
		config:          nil,
		manager:         manager,
	}
}
func (n *Node) AllNeighbours() []string {
	if n.IsRoot() {
		return n.children
	}
	return append(n.children, n.parent)
}

func (n *Node) UpdateGorumsConfig() {
	nodes := n.AllNeighbours()
	cfg, err := n.manager.NewConfiguration(&QSpec{numNodes: len(nodes)}, gorums.WithNodeList(nodes))
	if err != nil {
		log.Fatalln("failed to create gorums client config")
	}
	n.config = cfg

}

// SetNeighboursFromNodeList assumes a binary tree as slice where a nodes children are at index 2i+1 and 2i+2
func (n *Node) SetNeighboursFromNodeList(nodes []string) {
	for i := 0; i < len(nodes); i++ {
		// find n as a child of current node -> current node is n's parent
		if len(nodes) > (2*i+1) && nodes[2*i+1] == n.id {
			n.parent = nodes[i]
			continue
		}
		if len(nodes) > (2*i+2) && nodes[2*i+2] == n.id {
			n.parent = nodes[i]
			continue
		}

		// find n -> 2i+1 and 2i+2 are n's children if they exist
		if nodes[i] == n.id {
			if len(nodes) > (2*i + 1) {
				n.children = append(n.children, nodes[2*i+1])
			}
			if len(nodes) > (2*i + 2) {
				n.children = append(n.children, nodes[2*i+2])
			}
			continue
		}

	}
}
func (n *Node) IsRoot() bool {
	return n.parent == ""
}

// Run starts the main loop of the node
func (n *Node) Run() {
	n.startGorumsServer(n.rpcAddr)

	for {
		select {
		case <-time.After(time.Second * 1):
			log.Println("send heartbeat here...")
		}
	}
}

func (n *Node) startGorumsServer(addr string) {
	n.gorumsServer = gorums.NewServer()
	protos.RegisterKeyValueServiceServer(n.gorumsServer, n)
	listener, listenErr := net.Listen("tcp", addr)
	if listenErr != nil {
		log.Panicf("could not listen to address %v", addr)
	}
	go func() {
		serveErr := n.gorumsServer.Serve(listener)
		if serveErr != nil {
			log.Panicf("gorums server could not serve key value server")
		}
	}()
}

func (n *Node) Write(ctx gorums.ServerCtx, request *protos.WriteRequest) (response *emptypb.Empty, err error) {
	n.keyValueStorage.WriteValue(request.GetID(), request.GetKey(), request.GetValue())
	return &emptypb.Empty{}, nil
}

func (n *Node) Read(ctx gorums.ServerCtx, request *protos.ReadRequest) (response *protos.ReadResponse, err error) {
	value, err := n.keyValueStorage.ReadValue(request.Key)
	if err != nil {
		return &protos.ReadResponse{Value: 0}, err
	}
	return &protos.ReadResponse{Value: value}, nil
}

func (n *Node) ReadAll(ctx gorums.ServerCtx, request *protos.ReadRequest) (response *protos.ReadAllResponse, err error) {
	//TODO implement me
	panic("implement me")
}
