package oncetree

type Node struct {
	rpcAddr         string
	keyValueStorage *KeyValueStorage
}

func NewNode(rpcAddr string) *Node {
	return &Node{
		rpcAddr:         rpcAddr,
		keyValueStorage: NewKeyValueStorage(),
	}
}

func (n *Node) Run() {
}
