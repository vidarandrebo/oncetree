package oncetree

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"sync"
	"time"

	"github.com/relab/gorums"
	"github.com/vidarandrebo/oncetree/protos"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Node struct {
	rpcAddr         string
	keyValueStorage *KeyValueStorage
	gorumsServer    *gorums.Server
	id              string
	parent          map[string]string
	children        map[string]string
	gorumsManager   *protos.Manager
	gorumsConfig    *protos.Configuration
	logger          *log.Logger
	timestamp       int64
	mut             sync.Mutex
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
		gorumsConfig:    nil,
		gorumsManager:   manager,
		logger:          log.New(os.Stderr, fmt.Sprintf("NodeID: %s ", id), log.Ltime|log.Lmsgprefix),
		timestamp:       0,
	}
}

// Run starts the main loop of the node
func (n *Node) Run() {
	n.startGorumsServer(n.rpcAddr)

	for {
		select {
		case <-time.After(time.Second * 1):
			n.SendHeartbeat()
		}
	}
}

func (n *Node) allNeighbourAddrs() []string {
	neighbours := make([]string, 0)
	for _, address := range n.parent {
		neighbours = append(neighbours, address)
	}
	for _, address := range n.children {
		neighbours = append(neighbours, address)
	}
	return neighbours
}

func (n *Node) allNeighbourIDs() []string {
	IDs := make([]string, 0)
	for id := range n.parent {
		IDs = append(IDs, id)
	}
	for id := range n.children {
		IDs = append(IDs, id)
	}
	return IDs
}

// TryUpdateGorumsConfig will try to update the configuration 5 times with increasing timeout between attempts
//
// Non blocking fn
func (n *Node) TryUpdateGorumsConfig() {
	go func() {
		delays := []int64{1, 2, 3, 4, 5}
		for _, delay := range delays {
			nodes := n.allNeighbourAddrs()
			cfg, err := n.gorumsManager.NewConfiguration(&QSpec{numNodes: len(nodes)}, gorums.WithNodeList(nodes))
			if err == nil {
				n.logger.Println("created new gorums configuration")
				n.gorumsConfig = cfg
				break
			}
			if err != nil {
				n.logger.Printf("failed to create gorums gorumsConfig, retrying in %d seconds", delay)
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

func (n *Node) isRoot() bool {
	return len(n.parent) == 0
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
	n.mut.Lock()
	n.timestamp++
	n.keyValueStorage.WriteValue(n.id, request.GetKey(), request.GetValue(), n.timestamp)
	n.mut.Unlock()
	go func() {
		n.SendGossip(n.id, request.GetKey())
	}()
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

func (n *Node) PrintState(ctx gorums.ServerCtx, request *emptypb.Empty) (response *emptypb.Empty, err error) {
	n.logger.Println(n.keyValueStorage)
	return &emptypb.Empty{}, nil
}

func (n *Node) Gossip(ctx gorums.ServerCtx, request *protos.GossipMessage) (response *emptypb.Empty, err error) {
	n.logger.Printf("received gossip %v", request)
	n.mut.Lock()
	updated := n.keyValueStorage.WriteValue(request.GetNodeID(), request.GetKey(), request.GetValue(), request.GetTimestamp())
	n.mut.Unlock()
	if updated {
		go func() {
			n.SendGossip(request.NodeID, request.GetKey())
		}()
	}
	return &emptypb.Empty{}, nil
}

func (n *Node) Heartbeat(ctx gorums.ServerCtx, request *protos.HeartbeatMessage) {
	n.logger.Printf("received heartbeat from %s", request.GetNodeID())
}

func (n *Node) SendGossip(originID string, key int64) {
	n.mut.Lock()
	n.timestamp++
	ts := n.timestamp // make sure all messages has same ts
	n.mut.Unlock()
	for _, node := range n.gorumsConfig.Nodes() {
		nodeID, err := n.resolveNodeIDFromAddress(node.Address())
		if err != nil {
			continue
		}
		// skip returning to originID and sending to self.
		if nodeID == originID || nodeID == n.id {
			continue
		}
		value, err := n.keyValueStorage.ReadValueExceptNode(nodeID, key)
		// TODO handle err
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		_, err = node.Gossip(ctx, &protos.GossipMessage{NodeID: n.id, Key: key, Value: value, Timestamp: ts})
		// TODO handle err
		cancel()
	}
}

func (n *Node) SendHeartbeat() {
	go func() {
		if n.gorumsConfig == nil {
			return
		}
		msg := protos.HeartbeatMessage{NodeID: n.id}
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		n.gorumsConfig.Heartbeat(ctx, &msg)
	}()
}

func (n *Node) resolveNodeIDFromAddress(address string) (string, error) {
	for id, parentAddress := range n.parent {
		if parentAddress == address {
			return id, nil
		}
	}
	for id, childAddress := range n.children {
		if childAddress == address {
			return id, nil
		}
	}
	return "", fmt.Errorf("node with address %s not found", address)
}
