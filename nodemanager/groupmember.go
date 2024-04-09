package nodemanager

import "fmt"

type GroupMember struct {
	ID      string
	Address string
	Role    NodeRole
}

func (gm GroupMember) String() string {
	return fmt.Sprintf("GroupMember: { ID: %s, Address: %s, Role: %d }", gm.ID, gm.Address, gm.Role)
}

func NewGroupMember(id string, address string, role NodeRole) GroupMember {
	return GroupMember{
		ID:      id,
		Address: address,
		Role:    role,
	}
}
