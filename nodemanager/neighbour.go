package nodemanager

import "fmt"

type Neighbour struct {
	ID       string
	GorumsID uint32
	Address  string
	Group    Group
	Role     NodeRole
}

func NewNeighbour(ID string, gorumsID uint32, address string, role NodeRole) *Neighbour {
	return &Neighbour{
		ID:       ID,
		GorumsID: gorumsID,
		Address:  address,
		Group:    NewGroup(),
		Role:     role,
	}
}

func (n *Neighbour) String() string {
	return fmt.Sprintf("Neighbour: { ID: %s, GorumsID: %d, Address: %s, Role: %d }", n.ID, n.GorumsID, n.Address, n.Role)
}

func (n *Neighbour) GroupMemberIDs() []string {
	ids := make([]string, 0)
	for _, member := range n.Group.members {
		ids = append(ids, member.ID)
	}
	return ids
}
