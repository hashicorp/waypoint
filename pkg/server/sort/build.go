// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package sort

import (
	"sort"

	pb "github.com/hashicorp/waypoint/pkg/server/gen"
)

// BuildStartDesc sorts builds by start time descending (most recent first).
// For the opposite, use sort.Reverse.
type BuildStartDesc []*pb.Build

func (s BuildStartDesc) Len() int      { return len(s) }
func (s BuildStartDesc) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s BuildStartDesc) Less(i, j int) bool {
	t1 := s[i].Status.StartTime.AsTime()
	t2 := s[j].Status.StartTime.AsTime()

	return t2.Before(t1)
}

var (
	_ sort.Interface = (BuildStartDesc)(nil)
)
