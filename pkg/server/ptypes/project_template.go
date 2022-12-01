package ptypes

import (
	"github.com/imdario/mergo"
	"github.com/mitchellh/go-testing-interface"
	"github.com/stretchr/testify/require"

	pb "github.com/hashicorp/waypoint/pkg/server/gen"
)

// TestProjectTemplate returns a valid project template for tests.
func TestProjectTemplate(t testing.T, src *pb.ProjectTemplate) *pb.ProjectTemplate {
	t.Helper()

	if src == nil {
		src = &pb.ProjectTemplate{}
	}

	require.NoError(t, mergo.Merge(src, &pb.Project{
		Name: "test",
	}))

	return src
}

// TODO: validations
