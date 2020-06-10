package sort

import (
	"sort"

	pb "github.com/hashicorp/waypoint/internal/server/gen"
)

// ConfigName sorts config variables by name.
type ConfigName []*pb.ConfigVar

func (s ConfigName) Len() int      { return len(s) }
func (s ConfigName) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s ConfigName) Less(i, j int) bool {
	return s[i].Name < s[j].Name
}

var (
	_ sort.Interface = (ConfigName)(nil)
)
