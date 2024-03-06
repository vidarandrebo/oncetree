package nodemanager

type NodeRole int

const (
	Parent NodeRole = iota
	Child
	Tmp
)
