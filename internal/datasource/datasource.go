// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

// Package datasource has the packages and logic that Waypoint uses
// for sourcing data for remote runs.
package datasource

import (
	"context"
	"reflect"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/hcl/v2"

	"github.com/hashicorp/waypoint-plugin-sdk/terminal"
	pb "github.com/hashicorp/waypoint/pkg/server/gen"
)

// Sourcer is implemented by all data sourcers and is responsible for
// sourcing data, configuring projects, determining the default values
// for operations, and more.
type Sourcer interface {
	// ProjectSource translates the configuration into a default data
	// source for the project. This should also perform any validation
	// on the configuration.
	ProjectSource(hcl.Body, *hcl.EvalContext) (*pb.Job_DataSource, error)

	// Override reconfigures the given data source with the given overrides.
	Override(*pb.Job_DataSource, map[string]string) error

	// RefToOverride converts a ref to the override settings necessary
	// to use this ref with a job.
	RefToOverride(*pb.Job_DataSource_Ref) (map[string]string, error)

	// Get downloads the sourced data and returns the directory where
	// the data is stored, a cleanup function, and any errors that occurred.
	// The cleanup function may be nil.
	Get(
		ctx context.Context,
		log hclog.Logger,
		ui terminal.UI,
		source *pb.Job_DataSource,
		baseDir string,
	) (string, *pb.Job_DataSource_Ref, func() error, error)

	// Changes is called to check a source for changes. The "current" argument
	// contains the current ref in use by the default workspace. The return
	// value should be non-nil only if there is a new ref to set.
	//
	// The boolean return value is the "ignore" boolean: if this is true,
	// there IS a change detected in the exact ref, but it should be ignored.
	// In this scenario, it means that the ref change doesn't result in a
	// triggerable-change (i.e. maybe the files didn't change at all).
	//
	// NOTE(mitchellh): I'm fairly sure that this will change when we
	// introduce per-workspace refs and stuff but given this is all internal
	// we're going with the crufty-but-works approach first.
	Changes(
		ctx context.Context,
		log hclog.Logger,
		ui terminal.UI,
		source *pb.Job_DataSource,
		current *pb.Job_DataSource_Ref,
		tempDir string,
	) (*pb.Job_DataSource_Ref, bool, error)
}

var (
	// FromString maps a string key to a source implementation.
	FromString = map[string]func() Sourcer{
		"git":   newGitSource,
		"local": newLocalSource,
	}

	// FromType maps a server DataSource type to a source implementation.
	FromType = map[reflect.Type]func() Sourcer{
		reflect.TypeOf((*pb.Job_DataSource_Git)(nil)):   newGitSource,
		reflect.TypeOf((*pb.Job_DataSource_Local)(nil)): newLocalSource,
	}
)
