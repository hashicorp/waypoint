package ptypes

import (
	"github.com/imdario/mergo"
	"github.com/mitchellh/go-testing-interface"
	"github.com/stretchr/testify/require"

	pb "github.com/hashicorp/waypoint/internal/server/gen"
)

func TestJobNew(t testing.T, src *pb.Job) *pb.Job {
	t.Helper()

	if src == nil {
		src = &pb.Job{}
	}

	require.NoError(t, mergo.Merge(src, &pb.Job{
		Application: &pb.Ref_Application{
			Application: "a_test",
			Project:     "p_test",
		},
		TargetRunner: &pb.Ref_Runner{
			Target: &pb.Ref_Runner_Any{
				Any: &pb.Ref_RunnerAny{},
			},
		},
	}))

	return src
}
