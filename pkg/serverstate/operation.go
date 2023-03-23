// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package serverstate

import (
	"github.com/hashicorp/go-memdb"
	pb "github.com/hashicorp/waypoint/pkg/server/gen"
)

// ListOperationOptions are options that can be set for List calls on
// operations for filtering and limiting the response.
type ListOperationOptions struct {
	Application   *pb.Ref_Application
	Workspace     *pb.Ref_Workspace
	Status        []*pb.StatusFilter
	Order         *pb.OperationOrder
	PhysicalState pb.Operation_PhysicalState
	WatchSet      memdb.WatchSet
}

// ListOperationOption is an exported type to set configuration for listing operations.
type ListOperationOption func(opts *ListOperationOptions)

// ListWithStatusFilter sets a status filter.
func ListWithStatusFilter(f ...*pb.StatusFilter) ListOperationOption {
	return func(opts *ListOperationOptions) {
		opts.Status = f
	}
}

// ListWithOrder sets ordering on the list operation.
func ListWithOrder(f *pb.OperationOrder) ListOperationOption {
	return func(opts *ListOperationOptions) {
		opts.Order = f
	}
}

// ListWithPhysicalState sets ordering on the list operation.
func ListWithPhysicalState(f pb.Operation_PhysicalState) ListOperationOption {
	return func(opts *ListOperationOptions) {
		opts.PhysicalState = f
	}
}

// ListWithWorkspace sets ordering on the list operation.
func ListWithWorkspace(f *pb.Ref_Workspace) ListOperationOption {
	return func(opts *ListOperationOptions) {
		opts.Workspace = f
	}
}

// ListWithWatchSet registers watches for the listing, allowing the watcher
// to detect if new items are added.
func ListWithWatchSet(ws memdb.WatchSet) ListOperationOption {
	return func(opts *ListOperationOptions) {
		opts.WatchSet = ws
	}
}

// BuildListOperationOptions is a helper for implementations to create
// a ListOperationOptions from an app ref and a set of options.
func BuildListOperationOptions(ref *pb.Ref_Application, opts ...ListOperationOption) *ListOperationOptions {
	var result ListOperationOptions
	result.Application = ref
	for _, opt := range opts {
		opt(&result)
	}

	return &result
}
