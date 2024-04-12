package nodemanager

type NodeRole int

const (
	Parent NodeRole = iota
	Child
	Tmp      // Temporary node in join process
	Recovery // Part of the recovery group
)
