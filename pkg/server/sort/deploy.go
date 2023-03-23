// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package sort

import (
	"sort"

	pb "github.com/hashicorp/waypoint/pkg/server/gen"
)

// DeploymentStartDesc sorts deployments by start time descending.
type DeploymentStartDesc []*pb.Deployment

func (s DeploymentStartDesc) Len() int      { return len(s) }
func (s DeploymentStartDesc) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s DeploymentStartDesc) Less(i, j int) bool {
	t1 := s[i].Status.StartTime.AsTime()
	t2 := s[j].Status.StartTime.AsTime()

	return t2.Before(t1)
}

// DeploymentCompleteDesc sorts deployments by completion time descending.
type DeploymentCompleteDesc []*pb.Deployment

func (s DeploymentCompleteDesc) Len() int      { return len(s) }
func (s DeploymentCompleteDesc) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s DeploymentCompleteDesc) Less(i, j int) bool {
	t1 := s[i].Status.CompleteTime.AsTime()
	t2 := s[j].Status.CompleteTime.AsTime()

	return t2.Before(t1)
}

// DeploymentBundleStartDesc sorts deployment bundles by start time descending.
type DeploymentBundleStartDesc []*pb.UI_DeploymentBundle

func (s DeploymentBundleStartDesc) Len() int      { return len(s) }
func (s DeploymentBundleStartDesc) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s DeploymentBundleStartDesc) Less(i, j int) bool {
	t1 := s[i].Deployment.Status.StartTime.AsTime()
	t2 := s[j].Deployment.Status.StartTime.AsTime()

	return t2.Before(t1)
}

// DeploymentBundleCompleteDesc sorts deployment bundles by completion time descending.
type DeploymentBundleCompleteDesc []*pb.UI_DeploymentBundle

func (s DeploymentBundleCompleteDesc) Len() int      { return len(s) }
func (s DeploymentBundleCompleteDesc) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s DeploymentBundleCompleteDesc) Less(i, j int) bool {
	t1 := s[i].Deployment.Status.CompleteTime.AsTime()
	t2 := s[j].Deployment.Status.CompleteTime.AsTime()

	return t2.Before(t1)
}

var (
	_ sort.Interface = (DeploymentStartDesc)(nil)
	_ sort.Interface = (DeploymentCompleteDesc)(nil)
	_ sort.Interface = (DeploymentBundleStartDesc)(nil)
	_ sort.Interface = (DeploymentBundleCompleteDesc)(nil)
)
