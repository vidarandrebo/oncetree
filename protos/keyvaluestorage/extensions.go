package keyvaluestorage

func (c *Configuration) Node(id uint32) (*Node, bool) {
	for _, node := range c.nodes {
		if node.ID() == id {
			return node, true
		}
	}
	return nil, false
}
