package oncetree

import (
	"cmp"
	"fmt"
	"github.com/vidarandrebo/oncetree/concurrent/mutex"
	"log"
	"slices"
	"sync"

	"github.com/relab/gorums"
	nmprotos "github.com/vidarandrebo/oncetree/protos/nodemanager"
)

type NodeManager struct {
	id           string
	fanout       int
	neighbours   *ConcurrentMap[string, *Neighbour]
	epoch        int64
	nextGorumsID *mutex.RWMutex[uint32]
	lastJoinID   *mutex.RWMutex[string]
	joinMut      sync.Mutex
	logger       *log.Logger
}

func NewNodeManager(id string, fanout int) *NodeManager {
	return &NodeManager{
		id:           id,
		fanout:       fanout,
		neighbours:   NewConcurrentMap[string, *Neighbour](),
		nextGorumsID: mutex.New[uint32](0),
		lastJoinID:   mutex.New[string](""),
	}
}

func (nm *NodeManager) HandleFailure(nodeID string) {
}

func (nm *NodeManager) GetNeighbours() []KeyValuePair[string, *Neighbour] {
	return nm.neighbours.Entries()
}

func (nm *NodeManager) GetGorumsNeighbourMap() map[string]uint32 {
	IDs := make(map[string]uint32)
	for _, neighbour := range nm.neighbours.Values() {
		IDs[neighbour.Address] = neighbour.GorumsID
	}
	return IDs
}

func (nm *NodeManager) AddNeighbour(nodeID string, address string, role NodeRole) {
	nextID := nm.nextGorumsID.Lock()
	gorumsID := *nextID
	*nextID += 1
	nm.nextGorumsID.Unlock(&nextID)
	neighbour := NewNeighbour(nodeID, gorumsID, address, role)
	nm.setNeighbour(nodeID, neighbour)
}

func (nm *NodeManager) setNeighbour(nodeID string, neighbour *Neighbour) {
	nm.neighbours.Set(nodeID, neighbour)
}

func (nm *NodeManager) GetNeighbour(nodeID string) (*Neighbour, bool) {
	value, exists := nm.neighbours.Get(nodeID)
	return value, exists
}

func (nm *NodeManager) AllNeighbourIDs() []string {
	return nm.neighbours.Keys()
}

func (nm *NodeManager) AllNeighbourAddrs() []string {
	addresses := make([]string, 0)
	for _, neighbour := range nm.neighbours.Values() {
		addresses = append(addresses, neighbour.Address)
	}
	return addresses
}

func (nm *NodeManager) resolveNodeIDFromAddress(address string) (string, error) {
	for _, neighbour := range nm.neighbours.Entries() {
		if neighbour.Value.Address == address {
			return neighbour.Key, nil
		}
	}
	return "", fmt.Errorf("node with address %s not found", address)
}

func (nm *NodeManager) GetChildren() []*Neighbour {
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

func (nm *NodeManager) GetParent() *Neighbour {
	for _, neighbour := range nm.neighbours.Values() {
		if neighbour.Role == Parent {
			return neighbour
		}
	}
	return nil
}

func (nm *NodeManager) Join(ctx gorums.ServerCtx, request *nmprotos.JoinRequest) (response *nmprotos.JoinResponse, err error) {
	nm.joinMut.Lock()
	defer nm.joinMut.Unlock()
	response = &nmprotos.JoinResponse{
		OK:          false,
		NodeID:      nm.id,
		NextAddress: "",
	}
	children := nm.GetChildren()

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
	children := nm.GetChildren()
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
