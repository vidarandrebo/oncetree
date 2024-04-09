package nodemanager

import (
	"fmt"
	"strings"
)

type Group struct {
	epoch   int64
	members []GroupMember
}

func NewGroup() *Group {
	return &Group{
		epoch:   0,
		members: make([]GroupMember, 0),
	}
}

func (g *Group) MemberStrings() []string {
	members := make([]string, 0)
	for _, member := range g.members {
		members = append(members, member.String())
	}
	return members
}

func (g *Group) String() string {
	return fmt.Sprintf("Epoch: %d, members: { %v }", g.epoch, strings.Join(g.MemberStrings(), ", "))
}
