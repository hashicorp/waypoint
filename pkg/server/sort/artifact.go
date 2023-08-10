// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package sort

import (
	"sort"

	pb "github.com/hashicorp/waypoint/pkg/server/gen"
)

// ArtifactStartDesc sorts builds by start time descending (most recent first).
// For the opposite, use sort.Reverse.
type ArtifactStartDesc []*pb.PushedArtifact

func (s ArtifactStartDesc) Len() int      { return len(s) }
func (s ArtifactStartDesc) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s ArtifactStartDesc) Less(i, j int) bool {
	t1 := s[i].Status.StartTime.AsTime()
	t2 := s[j].Status.StartTime.AsTime()

	return t2.Before(t1)
}

var (
	_ sort.Interface = (ArtifactStartDesc)(nil)
)
