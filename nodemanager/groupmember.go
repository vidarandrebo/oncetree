package nodemanager

import (
	"fmt"

	"github.com/vidarandrebo/oncetree/nodemanager/nmenums"
)

type GroupMember struct {
	ID      string
	Address string
	Role    nmenums.NodeRole
}

func (gm GroupMember) String() string {
	return fmt.Sprintf("GroupMember: { ID: %s, Address: %s, Role: %d }", gm.ID, gm.Address, gm.Role)
}

func NewGroupMember(id string, address string, role nmenums.NodeRole) GroupMember {
	return GroupMember{
		ID:      id,
		Address: address,
		Role:    role,
	}
}
