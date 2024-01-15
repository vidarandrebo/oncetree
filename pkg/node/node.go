package node

import (
	"sync"
	"github.com/vidarandrebo/once-tree/pkg/oncetree"
	"github.com/vidarandrebo/once-tree/pkg/server"
)

type Node struct {
	Server *server.Server;
	OnceTree *oncetree.OnceTree;
}

func NewNode(httpAddr string, rpcAddr string) *Node {
	server := server.NewServer(httpAddr)
	oncetree := oncetree.OnceTree{}
	return &Node{Server: server, OnceTree: &oncetree}
}

func(n *Node) Run() {
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		n.Server.Run()
	}();

	wg.Wait()

}