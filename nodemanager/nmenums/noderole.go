package nmenums

type NodeRole int

func (nr NodeRole) String() string {
	switch nr {
	case Parent:
		return "Parent"
	case Child:
		return "Child"
	case Tmp:
		return "Tmp"
	case Recovery:
		return "Recovery"
	default:
		return "Unknown Role"
	}
}

const (
	Parent NodeRole = iota
	Child
	Tmp      // Temporary node in join process
	Recovery // Part of the recovery group
)
