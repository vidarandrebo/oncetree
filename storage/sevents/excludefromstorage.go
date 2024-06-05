package sevents

import "github.com/vidarandrebo/oncetree/common/hashset"

type ExcludeFromStorageEvent struct {
	ReadExcludeIDs []string
	GossipExcludes map[string]hashset.HashSet[string] // dest -> excludeID
}
