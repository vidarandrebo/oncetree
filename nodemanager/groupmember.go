package nodemanager

import "fmt"

type GroupMember struct {
	Address string
	Role    NodeRole
}

func (gm *GroupMember) String() string {
	return fmt.Sprintf("GroupMember: { Address: %s, Role: %d }", gm.Address, gm.Role)
}

func NewGroupMember(address string, role NodeRole) *GroupMember {
	return &GroupMember{
		Address: address,
		Role:    role,
	}
}
