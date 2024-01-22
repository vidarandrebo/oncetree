package oncetree

import (
	"fmt"
	"github.com/relab/gorums"
	"github.com/vidarandrebo/oncetree/protos"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/emptypb"
	"log"
	"net"
	"os"
	"time"
)

type Node struct {
	rpcAddr         string
	keyValueStorage *KeyValueStorage
	gorumsServer    *gorums.Server
	id              string
	parent          map[string]string
	children        map[string]string
	manager         *protos.Manager
	config          *protos.Configuration
	logger          *log.Logger
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
		parent:          make(map[string]string),
		children:        make(map[string]string),
		config:          nil,
		manager:         manager,
		logger:          log.New(os.Stderr, fmt.Sprintf("NodeID: %s ", id), log.Ltime|log.Lmsgprefix),
	}
}
func (n *Node) AllNeighbourAddrs() []string {
	neighbours := make([]string, 0)
	for _, address := range n.parent {
		neighbours = append(neighbours, address)
	}
	for _, address := range n.children {
		neighbours = append(neighbours, address)
	}
	return neighbours
}
func (n *Node) AllNeighbourIDs() []string {
	IDs := make([]string, 0)
	for id, _ := range n.parent {
		IDs = append(IDs, id)
	}
	for id, _ := range n.children {
		IDs = append(IDs, id)
	}
	return IDs
}

// TryUpdateGorumsConfig will try to update the configuration 5 times with increasing timeout between attempts
func (n *Node) TryUpdateGorumsConfig() {
	go func() {
		delays := []int64{1, 2, 3, 4, 5}
		for _, delay := range delays {
			nodes := n.AllNeighbourAddrs()
			cfg, err := n.manager.NewConfiguration(&QSpec{numNodes: len(nodes)}, gorums.WithNodeList(nodes))
			if err == nil {
				n.logger.Println("Created new gorums configuration")
				n.config = cfg
				break
			}
			if err != nil {
				n.logger.Printf("failed to create gorums config, retrying in %d seconds", delay)
				time.Sleep(time.Second * time.Duration(delay))
			}
		}
	}()

}

// SetNeighboursFromNodeMap assumes a binary tree as slice where a nodes children are at index 2i+1 and 2i+2
// Uses both a slice and a map to ensure consistent iteration order
func (n *Node) SetNeighboursFromNodeMap(nodeIDs []string, nodes map[string]string) {
	for i, nodeID := range nodeIDs {
		// find n as a child of current node -> current node is n's parent
		if len(nodeIDs) > (2*i+1) && nodeIDs[2*i+1] == n.id {
			n.parent[nodeID] = nodes[nodeID]
			continue
		}
		if len(nodeIDs) > (2*i+2) && nodeIDs[2*i+2] == n.id {
			n.parent[nodeID] = nodes[nodeID]
			continue
		}

		// find n -> 2i+1 and 2i+2 are n's children if they exist
		if nodeID == n.id {
			if len(nodeIDs) > (2*i + 1) {
				childId := nodeIDs[2*i+1]
				n.children[childId] = nodes[childId]
			}
			if len(nodeIDs) > (2*i + 2) {
				childId := nodeIDs[2*i+2]
				n.children[childId] = nodes[childId]
			}
			continue
		}

	}
	n.logger.Printf("parent: %v", n.parent)
	n.logger.Printf("children: %v", n.children)
}
func (n *Node) IsRoot() bool {
	return len(n.parent) == 0
}

// Run starts the main loop of the node
func (n *Node) Run() {
	n.startGorumsServer(n.rpcAddr)

	for {
		select {
		case <-time.After(time.Second * 1):
			//n.logger.Println("send heartbeat here...")
		}
	}
}

func (n *Node) startGorumsServer(addr string) {
	n.gorumsServer = gorums.NewServer()
	protos.RegisterKeyValueServiceServer(n.gorumsServer, n)
	listener, listenErr := net.Listen("tcp", addr)
	if listenErr != nil {
		n.logger.Panicf("could not listen to address %v", addr)
	}
	go func() {
		serveErr := n.gorumsServer.Serve(listener)
		if serveErr != nil {
			n.logger.Panicf("gorums server could not serve key value server")
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
	value, err := n.keyValueStorage.ReadValue(request.Key)
	if err != nil {
		return &protos.ReadAllResponse{Value: nil}, err
	}
	return &protos.ReadAllResponse{Value: map[string]int64{n.id: value}}, nil
}
