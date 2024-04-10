package nodemanager

import (
	"cmp"
	"context"
	"fmt"
	"log/slog"
	"reflect"
	"slices"
	"sync"

	"github.com/vidarandrebo/oncetree/gorumsprovider"

	"github.com/vidarandrebo/oncetree/consts"
	"github.com/vidarandrebo/oncetree/eventbus"

	"github.com/google/uuid"

	"github.com/vidarandrebo/oncetree/concurrent/maps"
	"github.com/vidarandrebo/oncetree/concurrent/mutex"

	"github.com/relab/gorums"
	nmprotos "github.com/vidarandrebo/oncetree/protos/nodemanager"
)

type NodeManager struct {
	id             string
	address        string
	fanout         int
	neighbours     *maps.ConcurrentMap[string, *Neighbour]
	epoch          *mutex.RWMutex[int64]
	nextGorumsID   *mutex.RWMutex[uint32]
	lastJoinID     *mutex.RWMutex[string]
	joinMut        sync.Mutex
	logger         *slog.Logger
	eventBus       *eventbus.EventBus
	gorumsProvider *gorumsprovider.GorumsProvider
}

func New(id string, address string, fanout int, logger *slog.Logger, eventBus *eventbus.EventBus, gorumsProvider *gorumsprovider.GorumsProvider) *NodeManager {
	nm := &NodeManager{
		id:             id,
		address:        address,
		fanout:         fanout,
		neighbours:     maps.NewConcurrentMap[string, *Neighbour](),
		nextGorumsID:   mutex.New[uint32](0),
		lastJoinID:     mutex.New(""),
		epoch:          mutex.New[int64](0),
		logger:         logger.With(slog.Group("node", slog.String("module", "nodemanager"))),
		eventBus:       eventBus,
		gorumsProvider: gorumsProvider,
	}
	eventBus.RegisterHandler(reflect.TypeOf(NeighbourReadyEvent{}),
		func(e any) {
			if event, ok := e.(NeighbourReadyEvent); ok {
				nm.HandleNeighbourReadyEvent(event)
			}
		})
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
}

func (nm *NodeManager) HandleNeighbourReadyEvent(e NeighbourReadyEvent) {
	// TODO - move this to separate function
	nm.logger.Info("added neighbour", "id", e.NodeID)
	epoch := nm.epoch.Lock()
	*epoch++
	epochVal := *epoch
	members := make([]*nmprotos.GroupMemberInfo, 0)
	for _, m := range nm.neighbours.Values() {
		members = append(members, &nmprotos.GroupMemberInfo{
			Role:    int64(m.Role),
			Address: m.Address,
			ID:      m.ID,
		})
	}

	nm.epoch.Unlock(&epoch)

	cfg := nm.gorumsProvider.NodeManagerConfig()
	ctx, cancel := context.WithTimeout(context.Background(), consts.RPCContextTimeout)
	defer cancel()
	cfg.GroupInfo(ctx, &nmprotos.GroupInfoMessage{
		Epoch:   epochVal,
		NodeID:  nm.id,
		Members: members,
	})
}

func (nm *NodeManager) HandleNeighbourRemovedEvent(e NeighbourRemovedEvent) {
	nm.logger.Info("removed neighbour", "id", e.NodeID)
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

// AddNeighbour creates and stores a new neighbour with the provided parameters.
// The node will be included in all configurations provided by
// the gorumsProvider after AddNeighbour returns.
// A NewNeigbourAddedEvent is pushed to eventbus if the node is not temporary
func (nm *NodeManager) AddNeighbour(nodeID string, address string, role NodeRole) {
	nextID := nm.nextGorumsID.Lock()
	gorumsID := *nextID
	*nextID += 1
	nm.nextGorumsID.Unlock(&nextID)
	neighbour := NewNeighbour(nodeID, gorumsID, address, role)
	nm.neighbours.Set(nodeID, neighbour)
	nm.gorumsProvider.SetNodes(nm.GorumsNeighbourMap())
	if role != Tmp {
		nm.eventBus.PushEvent(NewNeigbourAddedEvent(nodeID))
	}
}

func (nm *NodeManager) Neighbour(nodeID string) (*Neighbour, bool) {
	value, exists := nm.neighbours.Get(nodeID)
	return value, exists
}

func (nm *NodeManager) AllNeighbourIDs() []string {
	return nm.neighbours.Keys()
}

// TmpGorumsMap returns a map from address to gorumsID with all the temporary nodes stored
// This map should not have more than 1 entry
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

// clearTmp removes all nodes that has a temporary role
func (nm *NodeManager) clearTmp() {
	for _, node := range nm.neighbours.Values() {
		if node.Role == Tmp {
			nm.neighbours.Delete(node.ID)
		}
	}
}

func (nm *NodeManager) NeighbourAddresses() []string {
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
	return "", fmt.Errorf("[NodeManager] - node with address %s not found", address)
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

// Parent returns the Node's parent.
// Will return nil if the node is the root node, or is not part of a tree
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
		knownNodeID, _ := uuid.NewV7()
		nm.AddNeighbour(knownNodeID.String(), knownAddr, Tmp)
		ctx, cancel := context.WithTimeout(context.Background(), consts.RPCContextTimeout)
		nm.gorumsProvider.SetNodes(nm.TmpGorumsMap())
		cfg := nm.gorumsProvider.NodeManagerConfig()

		node := cfg.Nodes()[0]
		joinRequest := &nmprotos.JoinRequest{
			NodeID:  nm.id,
			Address: nm.address,
		}
		response, err := node.Join(ctx, joinRequest)
		if err != nil {
			nm.logger.Error("join failed",
				slog.Any("err", err),
				slog.String("address", knownAddr),
			)
			panic("join failed")
		}

		nm.clearTmp()
		nm.gorumsProvider.Reset()

		if response.OK {
			joined = true
			nm.AddNeighbour(response.NodeID, knownAddr, Parent)
			nm.SendReady(response.NodeID)
		} else {
			knownAddr = response.NextAddress
		}
		cancel()
	}
}

func (nm *NodeManager) SendReady(nodeID string) {
	cfg := nm.gorumsProvider.NodeManagerConfig()
	ctx, cancel := context.WithTimeout(context.Background(), consts.RPCContextTimeout)
	defer cancel()

	gorumsID, ok := nm.GorumsID(nodeID)
	if !ok {
		nm.logger.Error(
			"failed to lookup gorumsID",
			slog.String("id", nodeID),
		)
	}
	node, ok := cfg.Node(gorumsID)
	response, err := node.Ready(ctx, &nmprotos.ReadyMessage{NodeID: nm.id})
	if err != nil {
		nm.logger.Error(
			"failed to send ready message",
			slog.Any("err", err))
		return
	}
	if response.GetOK() {
		nm.eventBus.PushEvent(NewNeighbourReadyEvent(nodeID))
	}
}

// Join RPC will either accept the node as one of is children or return the address of another node.
// If max fanout is NOT reached, the node will respond OK and include its ID.
// If max fanout is reached, the node will respond NOT OK and include the address
// of one of its children. The node will alternate between its children in repeated calls to Join.
func (nm *NodeManager) Join(ctx gorums.ServerCtx, request *nmprotos.JoinRequest) (response *nmprotos.JoinResponse, err error) {
	nm.logger.Debug(
		"RPC Join",
		slog.String("id", request.GetNodeID()),
		slog.String("address", request.GetAddress()),
	)
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
		nm.logger.Error("cannot call NextJoinID when node has no children")
		panic("nextjoinId")
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

// Ready is used to signal to the parent that the newly joined node has created a gorums config and is ready to participate in the protocol. This might get handled a different way later...
func (nm *NodeManager) Ready(ctx gorums.ServerCtx, request *nmprotos.ReadyMessage) (response *nmprotos.ReadyMessage, err error) {
	nm.logger.Debug(
		"RPC Ready",
		slog.String("id", request.GetNodeID()),
	)
	// check node exists
	_, ok := nm.Neighbour(request.GetNodeID())
	if !ok {
		return &nmprotos.ReadyMessage{OK: false, NodeID: nm.id}, fmt.Errorf("node %s is not joined to this node", request.GetNodeID())
	}

	nm.eventBus.PushEvent(NewNeighbourReadyEvent(request.GetNodeID()))
	return &nmprotos.ReadyMessage{OK: true, NodeID: nm.id}, nil
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

func (nm *NodeManager) GroupInfo(ctx gorums.ServerCtx, request *nmprotos.GroupInfoMessage) {
	// message arrive one by one from same client in gorums, so should not need to lock for epoch compare
	node, ok := nm.neighbours.Get(request.GetNodeID())

	// Group exists and is up to date
	if ok && node.Group.epoch >= request.GetEpoch() {
		return
	}
	if !ok {
		nm.logger.Error("received group info from unknown node",
			slog.String("id", request.GetNodeID()))
		return
	}

	newMembers := make([]GroupMember, 0)
	for _, member := range request.GetMembers() {
		newMembers = append(newMembers, NewGroupMember(
			member.GetID(),
			member.GetAddress(),
			NodeRole(member.GetRole())),
		)
	}
	node.Group.epoch = request.GetEpoch()
	node.Group.members = newMembers
	nm.logger.Debug("group updated",
		slog.String("id", request.GetNodeID()),
		slog.String("group", node.Group.String()))
}
