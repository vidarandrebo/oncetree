package main

import (
	"github.com/vidarandrebo/oncetree"
)

var (
	nodeIDs = []string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9"}
	nodeMap = map[string]string{
		"0": ":9080",
		"1": ":9081",
		"2": ":9082",
		"3": ":9083",
		"4": ":9084",
		"5": ":9085",
		"6": ":9086",
		"7": ":9087",
		"8": ":9088",
		"9": ":9089",
	}
)

func main() {
	nodeAddrs := make([]string, 0)
	for _, nodeID := range nodeIDs {
		nodeAddrs = append(nodeAddrs, nodeMap[nodeID])
	}
	client := oncetree.NewClient(nodeAddrs)
	client.Run()
}
