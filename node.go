package oncetree

import (
	"context"
	"fmt"
	"github.com/relab/gorums"
	"github.com/vidarandrebo/oncetree/protos/failuredetectorprotos"
	"google.golang.org/protobuf/types/known/emptypb"
	"log"
	"net"
	"os"
	"sync"
)

type Node struct {
	id              string
	rpcAddr         string
	gorumsServer    *gorums.Server
	gorumsManagers  *GorumsManagers
	nodeManager     *NodeManager
	logger          *log.Logger
	timestamp       int64
	failureDetector *FailureDetector
	mut             sync.Mutex
	stopChan        chan string
	nodeFailureChan <-chan string
}

func NewNode(id string, rpcAddr string) *Node {
	logger := log.New(os.Stderr, fmt.Sprintf("NodeID: %s ", id), log.Ltime|log.Lmsgprefix)
	return &Node{
		rpcAddr:         rpcAddr,
		id:              id,
		failureDetector: NewFailureDetector(id, logger),
		nodeManager:     NewNodeManager(),
		gorumsManagers:  CreateGorumsManagers(),
		logger:          logger,
		timestamp:       0,
		stopChan:        make(chan string),
	}
}

func (n *Node) startGorumsServer(addr string) {
	n.gorumsServer = gorums.NewServer()
	failuredetectorprotos.RegisterFailureDetectorServiceServer(n.gorumsServer, n.failureDetector)
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

// Run starts the main loop of the node
func (n *Node) Run() {
	ctx, cancel := context.WithCancel(context.Background())
	wg := sync.WaitGroup{}
	defer cancel()
	n.startGorumsServer(n.rpcAddr)

	n.failureDetector.RegisterNodes(n.nodeManager.AllNeighbourIDs())
	n.nodeFailureChan = n.failureDetector.Subscribe()
	wg.Add(1)
	n.failureDetector.Run(ctx, &wg)

	nodeExitMessage := ""
mainLoop:
	for {
		select {
		case failedNode := <-n.nodeFailureChan:
			n.nodeManager.HandleFailure(failedNode)
		case nodeExitMessage = <-n.stopChan:
			n.logger.Println("Main loop stopped manually via stop channel")
			cancel()
			break mainLoop
		}
	}
	wg.Wait()
	n.gorumsServer.Stop()
	n.logger.Printf("Exiting after message \"%s\" on stop channel", nodeExitMessage)
}

// SetNeighboursFromNodeMap assumes a binary tree as slice where a nodes children are at index 2i+1 and 2i+2
// Uses both a slice and a map to ensure consistent iteration order
func (n *Node) SetNeighboursFromNodeMap(nodeIDs []string, nodes map[string]string) {
	for i, nodeID := range nodeIDs {
		// find n as a child of current node -> current node is n's parent
		if len(nodeIDs) > (2*i+1) && nodeIDs[2*i+1] == n.id {
			n.nodeManager.SetNeighbour(nodeID, NewNeighbour(nodes[nodeID], Parent))
			continue
		}
		if len(nodeIDs) > (2*i+2) && nodeIDs[2*i+2] == n.id {
			n.nodeManager.SetNeighbour(nodeID, NewNeighbour(nodes[nodeID], Parent))
			continue
		}

		// find n -> 2i+1 and 2i+2 are n's children if they exist
		if nodeID == n.id {
			if len(nodeIDs) > (2*i + 1) {
				childId := nodeIDs[2*i+1]
				n.nodeManager.SetNeighbour(childId, NewNeighbour(nodes[childId], Child))
			}
			if len(nodeIDs) > (2*i + 2) {
				childId := nodeIDs[2*i+2]
				n.nodeManager.SetNeighbour(childId, NewNeighbour(nodes[childId], Child))
			}
			continue
		}

	}
	n.logger.Printf("parent: %v", n.nodeManager.GetParent())
	n.logger.Printf("children: %v", n.nodeManager.GetChildren())
}

func (n *Node) isRoot() bool {
	return n.nodeManager.GetParent() == nil
}

func (n *Node) stop(msg string) {
	n.stopChan <- msg
}

type NodeRole int

const (
	Parent NodeRole = iota
	Child
)

type Neighbour struct {
	Address string
	Group   map[string]*GroupMember
	Role    NodeRole
}

func NewNeighbour(address string, role NodeRole) *Neighbour {
	return &Neighbour{
		Address: address,
		Group:   make(map[string]*GroupMember),
		Role:    role,
	}
}

type GroupMember struct {
	Address string
	Role    NodeRole
}

func NewGroupMember(address string, role NodeRole) *GroupMember {
	return &GroupMember{
		Address: address,
		Role:    role,
	}
}

func (n *Node) Crash(ctx gorums.ServerCtx, request *emptypb.Empty) (response *emptypb.Empty, err error) {
	n.stop("crash RPC")
	return &emptypb.Empty{}, nil
}
