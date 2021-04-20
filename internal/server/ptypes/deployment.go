package ptypes

import (
	"strconv"

	pb "github.com/hashicorp/waypoint/internal/server/gen"
)

// Type wrapper around the proto type so that we can add some methods.
type Deployment struct{ *pb.Deployment }

func (v *Deployment) URLFragment() string {
	// For older deployments (pre WP 0.4.0) we use the sequence. If
	// we have a generation set, we use the generation initial sequence.
	seq := v.Sequence

	// By using the generation sequence, we ensure that all deployments
	// in a generation share the same URL.
	if g := v.Generation; g != nil {
		seq = g.InitialSequence
	}

	return "v" + strconv.FormatUint(seq, 10)
}
