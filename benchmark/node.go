package benchmark

import "fmt"

type Node struct {
	ID       string
	Address  string
	GorumsID uint32
}

func (bn Node) String() string {
	return fmt.Sprintf("ID: %s, Address: %s, GorumsID %d", bn.ID, bn.Address, bn.GorumsID)
}

func GorumsMap(nodes map[string]Node) map[string]uint32 {
	gorumsMap := make(map[string]uint32)
	for _, node := range nodes {
		gorumsMap[node.Address] = node.GorumsID
	}
	return gorumsMap
}
