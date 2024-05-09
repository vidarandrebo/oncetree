package nodemanager

import (
	"cmp"
	"context"
	"fmt"
	"log/slog"
	"reflect"
	"slices"
	"sync"

	nodeprotos "github.com/vidarandrebo/oncetree/protos/node"

	"github.com/vidarandrebo/oncetree/concurrent/hashset"
	"github.com/vidarandrebo/oncetree/failuredetector/fdevents"
	"github.com/vidarandrebo/oncetree/nodemanager/nmenums"
	"github.com/vidarandrebo/oncetree/nodemanager/nmevents"

	"github.com/vidarandrebo/oncetree/gorumsprovider"

	"github.com/vidarandrebo/oncetree/eventbus"

	"github.com/google/uuid"

	"github.com/vidarandrebo/oncetree/concurrent/maps"
	"github.com/vidarandrebo/oncetree/concurrent/mutex"

	nmprotos "github.com/vidarandrebo/oncetree/protos/nodemanager"
)

type NodeManager struct {
	id              string
	address         string
	neighbours      *maps.ConcurrentMap[string, *Neighbour]
	epoch           *mutex.RWMutex[int64]
	nextGorumsID    *mutex.RWMutex[uint32]
	lastJoinID      *mutex.RWMutex[string]
	recoveryProcess *RecoveryProcess
	blackList       *hashset.ConcurrentHashSet[string]
	joinMut         sync.Mutex
	groupMut        sync.Mutex
	logger          *slog.Logger
	eventBus        *eventbus.EventBus
	gorumsProvider  *gorumsprovider.GorumsProvider
}

func New(id string, address string, logger *slog.Logger, eventBus *eventbus.EventBus, gorumsProvider *gorumsprovider.GorumsProvider) *NodeManager {
	nm := &NodeManager{
		id:              id,
		address:         address,
		neighbours:      maps.NewConcurrentMap[string, *Neighbour](),
		nextGorumsID:    mutex.New[uint32](0),
		lastJoinID:      mutex.New(""),
		epoch:           mutex.New[int64](0),
		logger:          logger.With(slog.Group("node", slog.String("module", "nodemanager"))),
		recoveryProcess: &RecoveryProcess{},
		blackList:       hashset.New[string](),
		eventBus:        eventBus,
		gorumsProvider:  gorumsProvider,
	}
	eventBus.RegisterHandler(reflect.TypeOf(nmevents.NeighbourAddedEvent{}),
		func(e any) {
			if event, ok := e.(nmevents.NeighbourAddedEvent); ok {
				nm.HandleNeighbourAddedEvent(event)
			}
		})
	eventBus.RegisterHandler(reflect.TypeOf(nmevents.NeighbourRemovedEvent{}),
		func(e any) {
			if event, ok := e.(nmevents.NeighbourRemovedEvent); ok {
				nm.HandleNeighbourRemovedEvent(event)
			}
		})

	eventBus.RegisterHandler(reflect.TypeOf(fdevents.NodeFailedEvent{}),
		func(e any) {
			if event, ok := e.(fdevents.NodeFailedEvent); ok {
				nm.logger.Info("handling node failure event",
					slog.String("id", event.NodeID))
				nm.SendPrepare(event)
			}
		})
	return nm
}

func (nm *NodeManager) SendPrepare(e fdevents.NodeFailedEvent) {
	nm.logger.Info("sending prepare")
	nm.recoveryProcess.mut.Lock()
	defer nm.recoveryProcess.mut.Unlock()

	if (nm.recoveryProcess.isActive) && (nm.recoveryProcess.groupID != e.NodeID) {
		panic("more than 1 concurrent failure, unrecoverable")
	}
	if nm.blackList.Contains(e.NodeID) {
		nm.logger.Info("node failure already handled, skipping",
			slog.String("id", e.NodeID))
		return
	}
	if !nm.recoveryProcess.isActive {
		nm.recoveryProcess.start(e.NodeID)
	}

	failedNode, ok := nm.neighbours.Get(e.NodeID)
	if !ok {
		nm.logger.Error("failed node does not exist",
			slog.String("id", e.NodeID))
		return
	}
	for _, groupMember := range failedNode.Group.members {
		if groupMember.ID == nm.id {
			nm.logger.Debug("group member is self, skipping")
			continue
		}
		_, neighbourExists := nm.Neighbour(groupMember.ID)
		if !neighbourExists {
			nm.AddNeighbour(groupMember.ID, groupMember.Address, nmenums.Recovery)
		}
	}
	recoveryMap := nm.GorumsRecoveryMap()
	if len(recoveryMap) == 0 {
		nm.logger.Info("no need to send prepare, only one node left in group")
		nm.neighbours.Delete(failedNode.ID)
		nm.eventBus.PushEvent(nmevents.NewNeigbourRemovedEvent(failedNode.ID))

		newGorumsNeighbourMap := nm.GorumsNeighbourMap()
		nm.gorumsProvider.ResetWithNewNodes(newGorumsNeighbourMap)

		nm.blackList.Add(failedNode.ID)
		nm.recoveryProcess.stop()
		nm.eventBus.PushTask(nm.SendGroupInfo)
		return
	}
	cfg, err := nm.gorumsProvider.CustomNodeManagerConfig(nm.GorumsRecoveryMap())
	if err != nil {
		nm.logger.Error("creation of recovery gorums config failed",
			slog.Any("err", err))
		panic("creation of recovery gorums config failed") // TODO - remove in prod
	}

	ctx := context.Background()
	response, err := cfg.Prepare(ctx, &nmprotos.PrepareMessage{
		NodeID:  nm.id,
		GroupID: failedNode.ID,
		Epoch:   failedNode.Group.epoch,
	})
	nm.logger.Info("received response from prepare",
		slog.Bool("isLeader", response.GetOK()))
	nm.recoveryProcess.isLeader = response.GetOK()

	if response.GetOK() {
		go nm.SendAccept()
		// nm.eventBus.PushTask(nm.SendAccept)
	}
}

func (nm *NodeManager) SendAccept() {
	nm.logger.Info("sending accept")

	nm.recoveryProcess.mut.Lock()
	defer nm.recoveryProcess.mut.Unlock()

	ctx := context.Background()
	cfg, err := nm.gorumsProvider.CustomNodeManagerConfig(nm.GorumsRecoveryMap())
	if err != nil {
		nm.logger.Error("creation of recovery gorums config failed",
			slog.Any("err", err))
		panic("creation of recovery gorums config failed") // TODO - remove in prod
	}
	node, ok := nm.neighbours.Get(nm.recoveryProcess.groupID)

	if !ok {
		panic("failed node not found")
	}

	// here we create the new tree structure, for the sake of simplicity, we just map all nodes to have self as parent,
	// but in more complex re-config scenario, we could nest the tree structure
	// we would then however make sure we don't introduce loops in the tree.
	newParent := make(map[string]string) // nodeID -> parentID
	for _, member := range node.GroupMemberIDs() {
		newParent[member] = nm.id
	}

	learn, err := cfg.Accept(ctx, &nmprotos.AcceptMessage{
		NodeID:    nm.id,
		GroupID:   nm.recoveryProcess.groupID,
		NewParent: newParent,
	})
	if err != nil {
		nm.logger.Error("Accept failed",
			slog.Any("err", err))
		panic("accept failed")
	}

	if learn.GetOK() {
		go nm.SendCommit()
		// nm.eventBus.PushTask(nm.SendCommit)
	}
}

func (nm *NodeManager) SendCommit() {
	recoveryMap := nm.GorumsRecoveryMap()
	cfg, err := nm.gorumsProvider.CustomNodeManagerConfig(recoveryMap)
	if err != nil {
		panic("could not create gorums config for recovery")
	}
	nm.recoveryProcess.mut.Lock()
	defer nm.recoveryProcess.mut.Unlock()
	ctx := context.Background()
	cfg.Commit(ctx, &nmprotos.CommitMessage{GroupID: nm.recoveryProcess.groupID})
	failedNode, ok := nm.neighbours.Get(nm.recoveryProcess.groupID)
	if !ok {
		panic("could not find failed node in neighbour map")
	}
	newChildren := make([]string, 0)
	for _, member := range failedNode.GroupMemberIDs() {
		if member == nm.id {
			continue
		}
		newChildren = append(newChildren, member)
	}
	for _, childID := range newChildren {
		newChild, ok := nm.neighbours.Get(childID)
		if !ok {
			panic("could not find new child in neighbour map")
		}
		newChild.Role = nmenums.Child
	}
	nm.neighbours.Delete(failedNode.ID)
	nm.eventBus.PushEvent(nmevents.NewNeigbourRemovedEvent(failedNode.ID))

	newGorumsNeighbourMap := nm.GorumsNeighbourMap()
	nm.gorumsProvider.ResetWithNewNodes(newGorumsNeighbourMap)

	for _, childID := range newChildren {
		nm.eventBus.PushEvent(nmevents.NewNeighbourAddedEvent(childID, "commit", nmenums.Child))
	}
	nm.blackList.Add(failedNode.ID)
	nm.recoveryProcess.stop()
	nm.eventBus.PushTask(nm.SendGroupInfo)
	nm.eventBus.PushEvent(nmevents.NewTreeRecoveredEvent(failedNode.ID))
}

func (nm *NodeManager) HandleNeighbourAddedEvent(e nmevents.NeighbourAddedEvent) {
	nm.logger.Info("added neighbour",
		slog.String("id", e.NodeID),
		slog.String("address", e.Address),
		slog.Any("role", e.Role))
}

func (nm *NodeManager) HandleNeighbourRemovedEvent(e nmevents.NeighbourRemovedEvent) {
	nm.logger.Info("removed neighbour", "id", e.NodeID)
}

func (nm *NodeManager) Neighbours() []maps.KeyValuePair[string, *Neighbour] {
	return nm.neighbours.Entries()
}

func (nm *NodeManager) NeighbourIDs() []string {
	IDs := make([]string, 0)
	for _, neighbour := range nm.neighbours.Values() {
		if (neighbour.Role == nmenums.Parent) || (neighbour.Role == nmenums.Child) {
			IDs = append(IDs, neighbour.ID)
		}
	}
	return IDs
}

func (nm *NodeManager) GorumsNeighbourMap() map[string]uint32 {
	IDs := make(map[string]uint32)
	for _, neighbour := range nm.neighbours.Values() {
		if (neighbour.Role == nmenums.Parent) || (neighbour.Role == nmenums.Child) {
			IDs[neighbour.Address] = neighbour.GorumsID
		}
	}
	return IDs
}

func (nm *NodeManager) GorumsRecoveryMap() map[string]uint32 {
	IDs := make(map[string]uint32)
	for _, neighbour := range nm.neighbours.Values() {
		if neighbour.Role == nmenums.Recovery {
			IDs[neighbour.Address] = neighbour.GorumsID
		}
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
// A NewNeighbourAddedEvent is pushed to eventbus if the node is not temporary
func (nm *NodeManager) AddNeighbour(nodeID string, address string, role nmenums.NodeRole) {
	nextID := nm.nextGorumsID.Lock()
	gorumsID := *nextID
	*nextID += 1
	nm.nextGorumsID.Unlock(&nextID)
	neighbour := NewNeighbour(nodeID, gorumsID, address, role)
	nm.neighbours.Set(nodeID, neighbour)
	nm.gorumsProvider.SetNodes(nm.GorumsNeighbourMap())
	if (role != nmenums.Tmp) && (role != nmenums.Recovery) {
		nm.eventBus.PushEvent(nmevents.NewNeighbourAddedEvent(nodeID, address, role))
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
		if node.Role == nmenums.Tmp {
			gorumsMap[node.Address] = node.GorumsID
			break
		}
	}
	return gorumsMap
}

// clearTmp removes all nodes that has a temporary role
func (nm *NodeManager) clearTmp() {
	for _, node := range nm.neighbours.Values() {
		if node.Role == nmenums.Tmp {
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
		if neighbour.Role == nmenums.Child {
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
		if neighbour.Role == nmenums.Parent {
			return neighbour
		}
	}
	return nil
}

func (nm *NodeManager) SendJoin(knownAddr string) {
	nm.joinMut.Lock()
	defer nm.joinMut.Unlock()
	if knownAddr == "" {
		return
	}
	joined := false
	for !joined {
		knownNodeID, _ := uuid.NewV7()
		nm.AddNeighbour(knownNodeID.String(), knownAddr, nmenums.Tmp)
		ctx := context.Background()
		nm.gorumsProvider.SetNodes(nm.TmpGorumsMap())
		cfg, ok := nm.gorumsProvider.NodeManagerConfig()
		if !ok {
			nm.logger.Error("failed to retrieve config for join operation", slog.String("fn", "nm.SendJoin"))
		}

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
			nm.AddNeighbour(response.GetNodeID(), response.GetAddress(), nmenums.Parent)
			nm.SendReady(response.GetNodeID())
		} else {
			knownAddr = response.NextAddress
		}
	}
}

func (nm *NodeManager) SendGroupInfo() {
	epoch := nm.epoch.Lock()
	*epoch++
	epochVal := *epoch
	members := make([]*nmprotos.GroupMemberInfo, 0)
	neighbours := nm.neighbours.Values()
	nm.epoch.Unlock(&epoch)
	for _, m := range neighbours {
		members = append(members, &nmprotos.GroupMemberInfo{
			Role:    int64(m.Role),
			Address: m.Address,
			ID:      m.ID,
		})
	}

	cfg, ok := nm.gorumsProvider.NodeManagerConfig()
	if !ok {
		nm.logger.Info("no nodes in config, skip sending GroupInfo")
		return
	}
	ctx := context.Background()
	cfg.GroupInfo(ctx, &nmprotos.GroupInfoMessage{
		Epoch:   epochVal,
		GroupID: nm.id,
		Members: members,
	})
}

func (nm *NodeManager) SendReady(nodeID string) {
	cfg, ok := nm.gorumsProvider.NodeManagerConfig()
	if !ok {
		nm.logger.Info("no nodes in config, skip sending Ready",
			slog.String("fn", "nm.SendReady"))
		panic("no nodes in config")
		return
	}
	ctx := context.Background()

	gorumsID, ok := nm.GorumsID(nodeID)
	if !ok {
		nm.logger.Error(
			"failed to lookup gorumsID",
			slog.String("id", nodeID),
		)
	}
	node, ok := cfg.Node(gorumsID)
	if !ok {
		panic("node not in config")
	}
	response, err := node.Ready(ctx, &nmprotos.ReadyMessage{NodeID: nm.id})
	if err != nil {
		nm.logger.Error(
			"failed to send ready message",
			slog.Any("err", err))
		return
	}
	if response.GetOK() {
		go nm.SendGroupInfo()
		nm.eventBus.PushEvent(nmevents.NewNeighbourReadyEvent(nodeID))
	}
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

func (nm *NodeManager) NodeConfig() (*nodeprotos.Configuration, bool) {
	return nm.gorumsProvider.NodeConfig()
}

func (nm *NodeManager) IsBlacklisted(nodeID string) bool {
	return nm.blackList.Contains(nodeID)
}
