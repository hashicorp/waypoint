package serverstate

import (
	pb "github.com/hashicorp/waypoint/internal/server/gen"
	"github.com/hashicorp/waypoint/internal/server/logbuffer"
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
