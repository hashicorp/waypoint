// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package sort

import (
	"sort"

	pb "github.com/hashicorp/waypoint/pkg/server/gen"
)

// ReleaseBundleCompleteDesc sorts deployment bundles by completion time descending.
type ReleaseBundleCompleteDesc []*pb.UI_ReleaseBundle

func (s ReleaseBundleCompleteDesc) Len() int      { return len(s) }
func (s ReleaseBundleCompleteDesc) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s ReleaseBundleCompleteDesc) Less(i, j int) bool {
	t1 := s[i].Release.Status.CompleteTime.AsTime()
	t2 := s[j].Release.Status.CompleteTime.AsTime()

	return t2.Before(t1)
}

var (
	_ sort.Interface = (ReleaseBundleCompleteDesc)(nil)
)
