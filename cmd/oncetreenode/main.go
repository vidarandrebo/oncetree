package main

import (
	"sync"

	"github.com/vidarandrebo/oncetree"
)

var (
	nodeIDs = []string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9", "10", "11", "12", "13", "14"}
	nodeMap = map[string]string{
		"0":  ":9080",
		"1":  ":9081",
		"2":  ":9082",
		"3":  ":9083",
		"4":  ":9084",
		"5":  ":9085",
		"6":  ":9086",
		"7":  ":9087",
		"8":  ":9088",
		"9":  ":9089",
		"10": ":9090",
		"11": ":9091",
		"12": ":9092",
		"13": ":9093",
		"14": ":9094",
	}
)

func main() {
	var wg sync.WaitGroup
	for _, id := range nodeIDs {
		id := id
		wg.Add(1)
		go func() {
			node := oncetree.NewNode(id, nodeMap[id])
			// node.SetNeighboursFromNodeMap(nodeIDs, nodeMap)
			if id == "0" {
				node.Run("")
			} else {
				node.Run(":9080")
			}
			wg.Done()
		}()
		// time.Sleep(1 * time.Second)
	}
	wg.Wait()
}
