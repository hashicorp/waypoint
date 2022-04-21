package ptypes

import (
	"github.com/imdario/mergo"
	"github.com/mitchellh/go-testing-interface"
	"github.com/stretchr/testify/require"

	pb "github.com/hashicorp/waypoint/pkg/server/gen"
)

// TestPipeline returns a valid user for tests.
func TestPipeline(t testing.T, src *pb.Pipeline) *pb.Pipeline {
	t.Helper()

	if src == nil {
		src = &pb.Pipeline{}
	}

	require.NoError(t, mergo.Merge(src, &pb.Pipeline{
		Id:   "test",
		Name: "test",
		Owner: &pb.Pipeline_Project{
			Project: &pb.Ref_Project{
				Project: "project",
			},
		},
		Steps: map[string]*pb.Pipeline_Step{
			"root": {
				Name: "root",
			},
		},
	}))

	return src
}
