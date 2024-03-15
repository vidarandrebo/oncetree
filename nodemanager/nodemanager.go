package nodemanager

import (
	"cmp"
	"context"
	"fmt"
	"log"
	"reflect"
	"slices"
	"sync"

	"github.com/vidarandrebo/oncetree/consts"
	"github.com/vidarandrebo/oncetree/eventbus"

	"github.com/google/uuid"

	"github.com/vidarandrebo/oncetree/concurrent/maps"
	"github.com/vidarandrebo/oncetree/concurrent/mutex"

	"github.com/relab/gorums"
	nmprotos "github.com/vidarandrebo/oncetree/protos/nodemanager"
)

type NodeManager struct {
	id            string
	address       string
	fanout        int
	neighbours    *maps.ConcurrentMap[string, *Neighbour]
	epoch         int64
	nextGorumsID  *mutex.RWMutex[uint32]
	lastJoinID    *mutex.RWMutex[string]
	joinMut       sync.Mutex
	logger        *log.Logger
	eventBus      *eventbus.EventBus
	gorumsManager *nmprotos.Manager
	gorumsConfig  *nmprotos.Configuration
}

func New(id string, address string, fanout int, logger *log.Logger, eventBus *eventbus.EventBus, gorumsManager *nmprotos.Manager) *NodeManager {
	nm := &NodeManager{
		id:            id,
		address:       address,
		fanout:        fanout,
		neighbours:    maps.NewConcurrentMap[string, *Neighbour](),
		nextGorumsID:  mutex.New[uint32](0),
		lastJoinID:    mutex.New(""),
		logger:        logger,
		eventBus:      eventBus,
		gorumsManager: gorumsManager,
	}
	eventBus.RegisterHandler(reflect.TypeOf(NeighbourAddedEvent{}),
		func(e any) {
			if event, ok := e.(NeighbourAddedEvent); ok {
				nm.HandleNeighbourAddedEvent(event)
			}
		})
	eventBus.RegisterHandler(reflect.TypeOf(NeighbourRemovedEvent{}),
		func(e any) {
			if event, ok := e.(NeighbourRemovedEvent); ok {
				nm.HandleNeighbourRemovedEvent(event)
			}
		})
	return nm
}

func (nm *NodeManager) HandleFailureEvent(nodeID string) {
}

func (nm *NodeManager) HandleNeighbourAddedEvent(e NeighbourAddedEvent) {
	nm.logger.Printf("added neighbour %s", e.NodeID)
}

func (nm *NodeManager) HandleNeighbourRemovedEvent(e NeighbourRemovedEvent) {
	nm.logger.Printf("removed neighbour %s", e.NodeID)
}

func (nm *NodeManager) Neighbours() []maps.KeyValuePair[string, *Neighbour] {
	return nm.neighbours.Entries()
}

func (nm *NodeManager) NeighbourIDs() []string {
	return nm.neighbours.Keys()
}

func (nm *NodeManager) GorumsNeighbourMap() map[string]uint32 {
	IDs := make(map[string]uint32)
	for _, neighbour := range nm.neighbours.Values() {
		IDs[neighbour.Address] = neighbour.GorumsID
	}
	return IDs
}

// GorumsID finds the gorumsID associated with the nodeID
func (nm *NodeManager) GorumsID(nodeID string) (uint32, bool) {
	node, ok := nm.neighbours.Get(nodeID)
	if ok {
		return node.GorumsID, true
	}
	return 0, false
}

// NodeID finds the nodeID associated with the gorumsID
func (nm *NodeManager) NodeID(gorumsID uint32) (string, bool) {
	for _, node := range nm.neighbours.Values() {
		if node.GorumsID == gorumsID {
			return node.ID, true
		}
	}
	return "", false
}

func (nm *NodeManager) AddNeighbour(nodeID string, address string, role NodeRole) {
	nextID := nm.nextGorumsID.Lock()
	gorumsID := *nextID
	*nextID += 1
	nm.nextGorumsID.Unlock(&nextID)
	neighbour := NewNeighbour(nodeID, gorumsID, address, role)
	nm.setNeighbour(nodeID, neighbour)
	if role != Tmp {
		nm.eventBus.Push(NewNeigbourAddedEvent(nodeID))
	}
}

func (nm *NodeManager) setNeighbour(nodeID string, neighbour *Neighbour) {
	nm.neighbours.Set(nodeID, neighbour)
}

func (nm *NodeManager) Neighbour(nodeID string) (*Neighbour, bool) {
	value, exists := nm.neighbours.Get(nodeID)
	return value, exists
}

func (nm *NodeManager) AllNeighbourIDs() []string {
	return nm.neighbours.Keys()
}

func (nm *NodeManager) TmpGorumsMap() map[string]uint32 {
	gorumsMap := make(map[string]uint32)
	for _, node := range nm.neighbours.Values() {
		if node.Role == Tmp {
			gorumsMap[node.Address] = node.GorumsID
			break
		}
	}
	return gorumsMap
}

func (nm *NodeManager) clearTmp() {
	for _, node := range nm.neighbours.Values() {
		if node.Role == Tmp {
			nm.neighbours.Delete(node.ID)
		}
	}
}

func (nm *NodeManager) AllNeighbourAddrs() []string {
	addresses := make([]string, 0)
	for _, neighbour := range nm.neighbours.Values() {
		addresses = append(addresses, neighbour.Address)
	}
	return addresses
}

func (nm *NodeManager) ResolveNodeIDFromAddress(address string) (string, error) {
	for _, neighbour := range nm.neighbours.Entries() {
		if neighbour.Value.Address == address {
			return neighbour.Key, nil
		}
	}
	return "", fmt.Errorf("node with address %s not found", address)
}

func (nm *NodeManager) Children() []*Neighbour {
	children := make([]*Neighbour, 0)
	for _, neighbour := range nm.neighbours.Values() {
		if neighbour.Role == Child {
			children = append(children, neighbour)
		}
	}
	slices.SortFunc(children, func(a, b *Neighbour) int {
		return cmp.Compare(a.ID, b.ID)
	})
	return children
}

func (nm *NodeManager) Parent() *Neighbour {
	for _, neighbour := range nm.neighbours.Values() {
		if neighbour.Role == Parent {
			return neighbour
		}
	}
	return nil
}

func (nm *NodeManager) SendJoin(knownAddr string) {
	if knownAddr == "" {
		return
	}
	joined := false
	for !joined {
		nm.gorumsManager.Close()
		knownNodeID, _ := uuid.NewV7()
		nm.AddNeighbour(knownNodeID.String(), knownAddr, Tmp)
		ctx, cancel := context.WithTimeout(context.Background(), consts.RPCContextTimeout)
		nodeMap := nm.TmpGorumsMap()
		cfg, err := nm.gorumsManager.NewConfiguration(
			&qspec{NumNodes: len(nodeMap)},
			gorums.WithNodeMap(nodeMap))
		if err != nil {
			nm.logger.Println("send join, create gorums config")
			nm.logger.Panicln(err)
		}
		nodes := cfg.Nodes()
		joinRequest := &nmprotos.JoinRequest{
			NodeID:  nm.id,
			Address: nm.address,
		}
		response, err := nodes[0].Join(ctx, joinRequest)
		if err != nil {
			nm.logger.Println(err)
			nm.logger.Fatalf("failed to join node with address %s", knownAddr)
		}
		if response.OK {
			joined = true
			nm.AddNeighbour(response.NodeID, knownAddr, Parent)
		} else {
			knownAddr = response.NextAddress
		}
		nm.clearTmp()
		cancel()
	}
}

func (nm *NodeManager) Join(ctx gorums.ServerCtx, request *nmprotos.JoinRequest) (response *nmprotos.JoinResponse, err error) {
	nm.joinMut.Lock()
	defer nm.joinMut.Unlock()
	response = &nmprotos.JoinResponse{
		OK:          false,
		NodeID:      nm.id,
		NextAddress: "",
	}
	children := nm.Children()

	// respond with one of the children's addresses if max fanout has been reached
	if len(children) == nm.fanout {
		nextPathID := nm.NextJoinID()
		nextPath, _ := nm.neighbours.Get(nextPathID)
		response.OK = false
		response.NextAddress = nextPath.Address
		return response, nil
	}
	response.OK = true
	response.NodeID = nm.id
	nm.AddNeighbour(request.NodeID, request.Address, Child)
	return response, nil
}

// NextJoinID returns the id of the node to send the next join request to
//
// Should only be called if max fanout has been reached, fn will panic if node has no children
func (nm *NodeManager) NextJoinID() string {
	lastJoinID := nm.lastJoinID.Lock()
	defer nm.lastJoinID.Unlock(&lastJoinID)
	children := nm.Children()
	if len(children) == 0 {
		nm.logger.Panicln("cannot call NextJoinID when node has no children")
	}
	lastJoinPathIndex := 0
	for i, child := range children {
		if child.ID == *lastJoinID {
			lastJoinPathIndex = i
		}
	}
	// last path was last child -> return first child
	if lastJoinPathIndex == len(children)-1 {
		*lastJoinID = children[0].ID
		return *lastJoinID
	}
	// increment last join path, then return
	*lastJoinID = children[lastJoinPathIndex+1].ID
	return *lastJoinID
}

func (nm *NodeManager) Prepare(ctx gorums.ServerCtx, request *nmprotos.PrepareMessage) (response *nmprotos.PromiseMessage, err error) {
	// TODO implement me
	panic("implement me")
}

func (nm *NodeManager) Accept(ctx gorums.ServerCtx, request *nmprotos.AcceptMessage) (response *nmprotos.LearnMessage, err error) {
	// TODO implement me
	panic("implement me")
}

func (nm *NodeManager) Commit(ctx gorums.ServerCtx, request *nmprotos.CommitMessage) {
	// TODO implement me
	panic("implement me")
}
