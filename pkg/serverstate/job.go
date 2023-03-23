// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package serverstate

import (
	"time"

	pb "github.com/hashicorp/waypoint/pkg/server/gen"
	"github.com/hashicorp/waypoint/pkg/server/logbuffer"
)

// These variables control the timeouts associated with the job system.
// The job system implementaions MUST use them. They are GUARANTEED
// to be written only when the state implementation is NOT running. Therefore,
// no lock is needed to read them.
//
// These MUST be used because the tests in statetest will manipulate these
// to verify various behaviors.
var (
	JobWaitingTimeout   = 2 * time.Minute
	JobHeartbeatTimeout = 2 * time.Minute
)

// Job is the exported structure that is returned for most state APIs
// and gives callers access to more information than the pure job structure.
type Job struct {
	// Full job structure.
	*pb.Job

	// OutputBuffer is the terminal output for this job. This is a buffer
	// that may not contain the full amount of output depending on the
	// time of connection.
	OutputBuffer *logbuffer.Buffer

	// Blocked is true if this job is blocked for some reason. The reasons
	// a job may be blocked:
	//  - another job for the same project/app/workspace.
	//  - a dependent job hasn't completed yet
	Blocked bool
}
