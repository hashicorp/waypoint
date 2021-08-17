package sort

import (
	"github.com/golang/protobuf/ptypes"
	pb "github.com/hashicorp/waypoint/internal/server/gen"
	"sort"
)

// ReleaseBundleCompleteDesc sorts deployment bundles by completion time descending.
type ReleaseBundleCompleteDesc []*pb.UI_ReleaseBundle

func (s ReleaseBundleCompleteDesc) Len() int      { return len(s) }
func (s ReleaseBundleCompleteDesc) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s ReleaseBundleCompleteDesc) Less(i, j int) bool {
	t1, err := ptypes.Timestamp(s[i].Release.Status.CompleteTime)
	if err != nil {
		return false
	}

	t2, err := ptypes.Timestamp(s[j].Release.Status.CompleteTime)
	if err != nil {
		return false
	}

	return t2.Before(t1)
}

var (
	_ sort.Interface = (ReleaseBundleCompleteDesc)(nil)
)
